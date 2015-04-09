package gohaveazurestorage

import "testing"

func TestRefactor(t *testing.T) {
	container := "containerfortesting"

	goHaveAzureStorage := NewWithDebug(Account, Key, true)
	blobStorage := goHaveAzureStorage.BlobStorage()

	httpStatusCode := blobStorage.CreateContainer(container)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	_, httpStatusCode = blobStorage.ListContainers()
	assertHTTPStatusCode(t, httpStatusCode, 200)

	_, httpStatusCode = blobStorage.GetContainerProperties(container)
	assertHTTPStatusCode(t, httpStatusCode, 200)

	_, httpStatusCode = blobStorage.GetContainerMetadata(container)
	assertHTTPStatusCode(t, httpStatusCode, 200)

	httpStatusCode = blobStorage.DeleteContainer(container)
	assertHTTPStatusCode(t, httpStatusCode, 202)
}
