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

package models

// ConfigFile represents a file containing config data needed to pull new Clam signature databases
// and config files from an S3 bucket.
type ConfigFile struct {
	// ClamMirrorBucket is the name of the bucket containing the clam DBs and configs.
	ClamMirrorBucket string `json:"clam_mirror_bucket"`
	// ClamConfigFiles is a simple implementation of a Set, so we can check the
	// bucket's file list against the list of files we actually want to download.
	ClamConfigFiles  []string `json:"clam_config_files"`
	ClamBucketKeyID  string   `json:"clam_bucket_key_id"`
	ClamBucketKey    string   `json:"clam_bucket_key"`
	ClamBucketRegion string   `json:"clam_bucket_region"`
}
