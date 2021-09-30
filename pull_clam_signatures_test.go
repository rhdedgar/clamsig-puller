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

package main_test

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rhdedgar/clamsig-puller"
	"github.com/rhdedgar/clamsig-puller/models"
)

type mockS3Client struct {
	s3iface.S3API
}

type mockDownloader struct {
	s3manager.Downloader
}

func (m *mockS3Client) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	curTime := time.Now()
	key0, key1, key2 := "test.cvd", "test2.cld", "test3.cvd"
	fSize := int64(1)

	return &s3.ListObjectsV2Output{
		Contents: []*s3.Object{{
			Key:          &key0,
			LastModified: &curTime,
			Size:         &fSize,
		}, {
			Key:          &key1,
			LastModified: &curTime,
			Size:         &fSize,
		}, {
			Key:          &key2,
			LastModified: &curTime,
			Size:         &fSize,
		},
		}}, nil
}

var _ = Describe("PullClamSignatures", func() {
	mockSvc := &mockS3Client{}
	mockConfig := &models.ConfigFile{
		ClamMirrorBucket: "testclammirrorbucket",
		ClamConfigFiles: []string{
			"file1",
			"file2",
			"file3",
		},
		ClamBucketKeyID:  "testclambucketkeyid",
		ClamBucketKey:    "testclambucketkey",
		ClamBucketRegion: "us-east-1",
	}

	Describe("GetSession", func() {
		Context("Validate the ability to generate an s3 session from a mock config", func() {
			It("Should use a config file to return a new AWS session", func() {
				sess, err := GetSession(mockConfig)

				Expect(err).To(BeNil())
				Expect(*sess.Config.Region).To(Equal("us-east-1"))
			})
		})
	})

	Describe("ListBucketObjects", func() {
		Context("Validate the ability to list files in an S3 bucket", func() {
			It("Should use a mock AWS session to get a slice of mock files", func() {
				v2OutPut, err := ListBucketObjects(mockSvc)

				Expect(err).To(BeNil())
				Expect(len(v2OutPut.Contents)).To(Equal(3))
			})
		})
	})

	Describe("GetService", func() {
		Context("Validate the ability to create a new s3.S3 client service", func() {
			It("Should use an existing session to return a new AWS client service", func() {
				sess, _ := GetSession(mockConfig)
				svc := GetService(sess)
				newSvc := s3.New(sess)

				Expect(*svc.Config.Region).To(Equal(*newSvc.Config.Region))
			})
		})
	})

	Describe("DownloadSignatures", func() {
		Context("Validate the ability to pull clam signature databases from an S3 bucket", func() {
			It("Should go through the list of mock bucket files and compare them with the file whitelist", func() {
				resp, _ := ListBucketObjects(mockSvc)

				err := DownloadSignatures(mockSvc, resp)
				if err != nil {
					fmt.Println("Error returned from DownloadSignatures:", err)
				}

				Expect(err).To(BeNil())
			})
		})
	})
})
