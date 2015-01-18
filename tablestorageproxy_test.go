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

  json1, _ := json.Marshal(entity)
  tableStorageProxy.InsertEntity(Table, json1)

  entity.RowKey = "456"
  json2, _ := json.Marshal(entity)
  tableStorageProxy.InsertEntity(Table, json2)

  entity.RowKey = "789"
  json3, _ := json.Marshal(entity)
  tableStorageProxy.InsertEntity(Table, json3)
}

func TestQueryEntity(t *testing.T) {
  goHaveStorage := New(Account, Key)
  tableStorageProxy := goHaveStorage.NewTableStorageProxy()

  tableStorageProxy.QueryEntity(Table, "ABC", "123", "")
}

func TestQueryEntityWithSelect(t *testing.T) {
  goHaveStorage := New(Account, Key)
  tableStorageProxy := goHaveStorage.NewTableStorageProxy()

  tableStorageProxy.QueryEntity(Table, "ABC", "123", "RowKey,Property1,Property3")
}

func TestQueryEntities(t *testing.T) {
  goHaveStorage := New(Account, Key)
  tableStorageProxy := goHaveStorage.NewTableStorageProxy()

  tableStorageProxy.QueryEntities(Table, "", "")
}

func TestQueryEntitiesWithSelect(t *testing.T) {
  goHaveStorage := New(Account, Key)
  tableStorageProxy := goHaveStorage.NewTableStorageProxy()

  tableStorageProxy.QueryEntities(Table, "RowKey,Property1,Property3", "")
}

func TestQueryEntitiesWithTop(t *testing.T) {
  goHaveStorage := New(Account, Key)
  tableStorageProxy := goHaveStorage.NewTableStorageProxy()

  tableStorageProxy.QueryEntities(Table, "", "1")
}

func TestQueryEntitiesWithSelectAndTop(t *testing.T) {
  goHaveStorage := New(Account, Key)
  tableStorageProxy := goHaveStorage.NewTableStorageProxy()

  tableStorageProxy.QueryEntities(Table, "RowKey,Property1,Property3", "1")
}

func DeleteTable(t *testing.T) {
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
