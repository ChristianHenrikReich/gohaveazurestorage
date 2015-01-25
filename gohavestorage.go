package gohavestorage

import (
	"encoding/base64"
	"gohavestorage/tablestorageproxy"
)

type GoHaveStorage struct {
	account      string
	key          []byte
	dumpSessions bool
}

func New(account string, key string) *GoHaveStorage {
	return NewWithDebug(account, key, false)
}

func NewWithDebug(account string, key string, dumpSessions bool) *GoHaveStorage {
	var goHaveStorage GoHaveStorage

	goHaveStorage.account = account
	goHaveStorage.key, _ = base64.StdEncoding.DecodeString(key)
	goHaveStorage.dumpSessions = dumpSessions

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

func (goHaveStorage *GoHaveStorage) DumpSessions() bool {
	return goHaveStorage.dumpSessions
}
