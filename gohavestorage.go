package gohavestorage

import (
	"encoding/base64"
	"gohavestorage/tablestorageproxy"
)

type GoHaveStorage struct {
	account string
	key     []byte
}

func New(account string, key string) *GoHaveStorage {
	var goHaveStorage GoHaveStorage

	goHaveStorage.account = account
	goHaveStorage.key, _ = base64.StdEncoding.DecodeString(key)

	return &goHaveStorage
}

func (goHaveStorage *GoHaveStorage) NewTableStorageProxy() *tablestorageproxy.TableStorageProxy {
	tableStorageProxy := tablestorageproxy.New(tablestorageproxy.GoHaveStorage(goHaveStorage))
	return tableStorageProxy
}

func (goHaveStorage *GoHaveStorage) GetAccount() string {
	return goHaveStorage.account
}

func (goHaveStorage *GoHaveStorage) GetKey() []byte {
	return goHaveStorage.key
}
