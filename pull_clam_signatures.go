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
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rhdedgar/clamsig-puller/config"
	"github.com/rhdedgar/clamsig-puller/models"
)

var (
	configFile = &config.ConfigFile
)

// GetSession returns an AWS S3 session.
func GetSession(configFile *models.ConfigFile) (*session.Session, error) {
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
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(configFile.ClamMirrorBucket)})
	if err != nil {
		fmt.Printf("Unable to list items in bucket %q", configFile.ClamMirrorBucket)
		return &s3.ListObjectsV2Output{}, err
	}

	return resp, nil
}

// DownloadSignatures compares signature databases on disk with those in the clam mirror bucket.
// It will download copies of the databases if found to be newer than what's on disk.
func DownloadSignatures(svc s3iface.S3API, resp *s3.ListObjectsV2Output) error {
	downloader := s3manager.NewDownloaderWithClient(svc)

	// Loop through bucket contents, and compare with our json array. If file is a match, then
	// check if doesn't exist, and check if the bucket's file is newer. Download it in those cases.
	for _, item := range resp.Contents {
		for _, localItem := range configFile.ClamConfigFiles {
			if *item.Key == localItem {
				fileStat, err := os.Stat(filepath.Join(config.ClamInstallDir, localItem))
				if os.IsNotExist(err) || fileStat.ModTime().Before(*item.LastModified) {
					// Adding a little jitter, so that there's not such a thundering herd
					// of network traffic upon updating.
					rand.Seed(time.Now().UnixNano())
					n := rand.Intn(5)
					time.Sleep(time.Duration(n) * time.Second)

					newFile, err := os.Create(filepath.Join(config.ClamInstallDir, localItem))
					if err != nil {
						fmt.Println("Unable to open file:", item)
						return err
					}

					defer newFile.Close()

					_, err = downloader.Download(newFile,
						&s3.GetObjectInput{
							Bucket: aws.String(configFile.ClamMirrorBucket),
							Key:    aws.String(*item.Key),
						})
					if err != nil {
						fmt.Println("Unable to download item:", item)
						return err
					}

					fmt.Println("Downloaded the following:")
					fmt.Println("Name:         ", *item.Key)
					fmt.Println("Last modified:", *item.LastModified)
					fmt.Println("Size:         ", *item.Size, "bytes")
					fmt.Println("")

				} else if err != nil {
					fmt.Println("Hit an issue opening the file:")
					return err
				}
			}
		}
	}
	return nil
}

func main() {
	fmt.Println("ClamAV signature and config updater v0.0.4")

	sess, err := GetSession(configFile)
	if err != nil {
		fmt.Println("Error returned from GetSession:", err)
	}

	svc := GetService(sess)

	resp, err := ListBucketObjects(svc)
	if err != nil {
		fmt.Println("Error returned from ListBucketObjects:", err)
	}

	err = DownloadSignatures(svc, resp)

	if err != nil {
		fmt.Println("Error returned from DownloadSignatures:", err)
	}
}
