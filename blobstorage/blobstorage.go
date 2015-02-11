package blobstorage

import "gohaveazurestorage/gohaveazurestoragecommon"

type BlobStorage struct {
	http             *gohaveazurestoragecommon.HTTP
	baseUrl          string
	secondaryBaseUrl string
}

func New(http *gohaveazurestoragecommon.HTTP) *BlobStorage {
	return &BlobStorage{http: http}
}
