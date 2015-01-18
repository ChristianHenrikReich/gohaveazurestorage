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

func InsertEntity(t *testing.T) {
  goHaveStorage := New(Account, Key)
  tableStorageProxy := goHaveStorage.NewTableStorageProxy()

  entity := &TestEntity{}
    entity.PartitionKey = "ABC"
    entity.RowKey = "123"
    entity.Property1 = "Value1"
    entity.Property2 = "Value2"
    entity.Property3 = "Value3"

    json, _ := json.Marshal(entity)

    tableStorageProxy.InsertEntity(Table, json)
  }

func TestQueryEntity(t *testing.T) {
  goHaveStorage := New(Account, Key)
  tableStorageProxy := goHaveStorage.NewTableStorageProxy()

  tableStorageProxy.QueryEntity(Table, "ABC", "123", "")
}

func TestQueryEntityWithSelect(t *testing.T) {
  goHaveStorage := New(Account, Key)
  tableStorageProxy := goHaveStorage.NewTableStorageProxy()

  tableStorageProxy.QueryEntity(Table, "ABC", "123", "Property1,Property3")
}

func TestDeleteTable(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()
	tableStorageProxy.DeleteTable(Table)
}

type TestEntity struct {
  PartitionKey string
  RowKey string
  Property1 string
  Property2 string
  Property3 string
}
