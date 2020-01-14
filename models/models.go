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
