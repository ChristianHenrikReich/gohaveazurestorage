package gohavestorage

import "testing"

var Key = ""
var Account = ""
var Table = "TestTable"

func CreateTable(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()
	tableStorageProxy.CreateTable(Table)
}

func QueryTables(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()
	tableStorageProxy.QueryTables()
}

func DeleteTable(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()
	tableStorageProxy.DeleteTable(Table)
}
