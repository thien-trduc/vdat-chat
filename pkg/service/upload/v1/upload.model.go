package upload

import (
	"github.com/minio/minio-go/v7"
	"time"
)

var FileGlobal = make(map[string]ModelUpload)
var minioClient *minio.Client

type ModelUpload struct {
	ShareUrl string
	Location string
	NameFile string
	Type     string
	CreateAt time.Time
}

type Info struct {
	Bucket           string
	Key              string
	ETag             string
	Size             int64
	LastModified     time.Time
	Location         string
	VersionID        string
	Expiration       time.Time
	ExpirationRuleID string
}
