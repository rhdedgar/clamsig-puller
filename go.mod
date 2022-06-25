module github.com/rhdedgar/clamsig-puller

go 1.16

replace github.com/rhdedgar/clamsig-puller/models => ./models

replace github.com/rhdedgar/clamsig-puller/datastores => ./datastores

require github.com/aws/aws-sdk-go v1.44.42
