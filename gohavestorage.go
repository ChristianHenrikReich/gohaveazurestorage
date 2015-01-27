package gohavestorage

import (
	"encoding/base64"
	"gohavestorage/gohavestoragecommon"
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
	http := gohavestoragecommon.NewHTTP("table", goHaveStorage.account, goHaveStorage.key, goHaveStorage.dumpSessions)
	return tablestorageproxy.New(http)
}
