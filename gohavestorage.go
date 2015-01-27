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
	decodedKey, _ := base64.StdEncoding.DecodeString(key)
	return &GoHaveStorage{account: account, key: decodedKey, dumpSessions: dumpSessions}
}

func (goHaveStorage *GoHaveStorage) NewTableStorageProxy() *tablestorageproxy.TableStorageProxy {
	http := gohavestoragecommon.NewHTTP("table", goHaveStorage.account, goHaveStorage.key, goHaveStorage.dumpSessions)
	return tablestorageproxy.New(http)
}
