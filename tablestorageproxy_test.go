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

func TestInsertEntity(t *testing.T) {
  goHaveStorage := New(Account, Key)
  tableStorageProxy := goHaveStorage.NewTableStorageProxy()

  entity := &TestEntity{}
  entity.PartitionKey = "ABC"
  entity.RowKey = "123"
  entity.Property1 = "Value1"

  json, _ := json.Marshal(entity)

  tableStorageProxy.InsertEntity(Table, json)
}

type TestEntity struct {
  PartitionKey string
  RowKey string
  Property1 string
}
