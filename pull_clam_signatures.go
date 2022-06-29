/*
Copyright 2020 Doug Edgar.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rhdedgar/clamsig-puller/datastores"
	"github.com/rhdedgar/clamsig-puller/models"
)

// GetSession returns an AWS S3 session.
func GetSession(configFile *models.AppSecrets) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(configFile.ClamBucketRegion),
		Credentials: credentials.NewStaticCredentials(configFile.ClamBucketKeyID, configFile.ClamBucketKey, ""),
	})
	if err != nil {
		fmt.Println("Error getting session:", err)
		return &session.Session{}, err
	}

	return sess, nil
}

// GetService returns a new S3 client service from an existing session.
func GetService(sess *session.Session) *s3.S3 {
	svc := s3.New(sess)

	return svc
}

// ListBucketObjects returns a list of an AWS s3 bucket's objects.
func ListBucketObjects(svc s3iface.S3API) (*s3.ListObjectsV2Output, error) {
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(datastores.AppSecrets.ClamMirrorBucket)})
	if err != nil {
		fmt.Printf("Unable to list items in bucket %q", datastores.AppSecrets.ClamMirrorBucket)
		return &s3.ListObjectsV2Output{}, err
	}

	return resp, nil
}

func dlToDisk(item string, downloader *s3manager.Downloader) (string, error) {
	newFile, err := os.Create(filepath.Join(datastores.AppSecrets.ClamConfigDir, item))
	if err != nil {
		fmt.Println("Unable to open file:", item)
		return "", err
	}
	defer newFile.Close()

	buf := aws.NewWriteAtBuffer([]byte{})

	_, err = downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(datastores.AppSecrets.ClamMirrorBucket),
			Key:    aws.String(item),
		})
	if err != nil {
		fmt.Println("Unable to download item:", item)
		return "", err
	}

	r := bytes.NewReader(buf.Bytes())

	zr, err := gzip.NewReader(r)
	if err != nil {
		fmt.Println(err)
	}

	sha := sha256.New()
	data := io.TeeReader(zr, sha)

	if _, err := io.Copy(newFile, data); err != nil {
		return "", fmt.Errorf("could not copy zr to newFile: %v\n", err)
	}

	if err := zr.Close(); err != nil {
		return "", fmt.Errorf("could not close zr: %v\n", err)
	}

	if err := newFile.Close(); err != nil {
		return "", fmt.Errorf("could not close newFile: %v\n", err)
	}

	checksum := hex.EncodeToString(sha.Sum(nil))

	return checksum, nil
}

func dlChecksumToString(item string, downloader *s3manager.Downloader) (string, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(datastores.AppSecrets.ClamMirrorBucket),
			Key:    aws.String(item),
		})
	if err != nil {
		fmt.Println("Unable to download item:", item)
		return "", err
	}

	return string(buf.Bytes()), nil
}

func loadConfigFile(filePath string, dest interface{}) error {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error loading secrets json from:  %v %v\n", filePath, err)
	}

	err = json.Unmarshal(fileBytes, dest)
	if err != nil {
		return fmt.Errorf("Error Unmarshaling secrets json: %v\n", err)
	}
	return nil
}

// DownloadSignatures compares signature databases on disk with those in the clam mirror bucket.
// It will download copies of the databases if found to be newer than what's on disk.
func DownloadSignatures(svc s3iface.S3API, resp *s3.ListObjectsV2Output) error {
	downloader := s3manager.NewDownloaderWithClient(svc)

	// Loop through bucket contents, and compare with our json array. If file is a match, then
	// check if doesn't exist, and check if the bucket's file is newer. Download it in those cases.
	for _, item := range resp.Contents {
		if strings.HasSuffix(*item.Key, ".gz") {
			for i := 0; i < 5; {
				splitItem := strings.Split(*item.Key, ".gz")
				baseItem := splitItem[0]

				// Adding a little jitter, so that there's not such a thundering herd
				// of network traffic upon updating.
				rand.Seed(time.Now().UnixNano())
				n := rand.Intn(i + 1)
				time.Sleep(time.Duration(n*10) * time.Second)

				bChecksum, err := dlChecksumToString(baseItem+"_checksum.txt", downloader)
				if err != nil {
					fmt.Printf("Hit an issue downloading checksum file: %v\n", err)
					fmt.Println("Proceeding without checksum file.")
				}

				fChecksum, err := dlToDisk(*item.Key, downloader)
				if err != nil {
					fmt.Printf("Skipping due to an issue downloading file: %v\n", err)
					continue
				}

				if bChecksum == "" && bChecksum != fChecksum {
					fmt.Println("Checksum mismatch; Possibly corrupted database file, trying again.")
				} else {
					fmt.Println("Downloaded the following:")
					fmt.Println("Name:         ", *item.Key)
					fmt.Println("Last modified:", *item.LastModified)
					fmt.Println("Size:         ", *item.Size, "bytes")
					fmt.Println("Continuing")
					i = 5
				}
			}
		}
	}
	return nil
}

func main() {
	fmt.Println("ClamAV signature and config updater v0.0.5")

	sess, err := GetSession(&datastores.AppSecrets)
	if err != nil {
		fmt.Println("Error returned from GetSession:", err)
		os.Exit(1)
	}

	svc := GetService(sess)

	resp, err := ListBucketObjects(svc)
	if err != nil {
		fmt.Println("Error returned from ListBucketObjects:", err)
		os.Exit(1)
	}

	err = DownloadSignatures(svc, resp)
	if err != nil {
		fmt.Println("Error returned from DownloadSignatures:", err)
		os.Exit(1)
	}

	fmt.Println("Finished running DownloadSignatures function.")
}
