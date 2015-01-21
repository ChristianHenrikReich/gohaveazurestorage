package gohavestorage

import (
	"encoding/json"
	"testing"
)

var Key = ""
var Account = ""
var Table = "TestTable"

func TestCreateTable(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()
	tableStorageProxy.CreateTable(Table)
}

func TestQueryTables(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()
	tableStorageProxy.QueryTables()
}

func TestInsertEntity(t *testing.T) {
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

	tableStorageProxy.QueryEntities(Table, "", "", "")
}

func TestQueryEntitiesWithSelect(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	tableStorageProxy.QueryEntities(Table, "RowKey,Property1,Property3", "", "")
}

func TestQueryEntitiesWithTop(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	tableStorageProxy.QueryEntities(Table, "", "", "1")
}

func TestQueryEntitiesWithSelectAndTop(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	tableStorageProxy.QueryEntities(Table, "RowKey,Property1,Property3", "", "1")
}

func TestQueryEntitiesWithSelectAndFilterAndTop(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	tableStorageProxy.QueryEntities(Table, "RowKey,Property1,Property3", "RowKey%20gt%20'123'", "1")
}

func TestDeleteEntity(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	tableStorageProxy.DeleteEntity(Table, "ABC", "123")
}

func TestUpdateEntity(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	entity := &TestEntity{}
	entity.PartitionKey = "ABC"
	entity.RowKey = "123"
	entity.Property1 = "Value1"
	entity.Property2 = "Value2"
	entity.Property3 = "Value3"

	json, _ := json.Marshal(entity)
	tableStorageProxy.UpdateEntity(Table, "ABC", "456", json)
}

func TestMergeEntity(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	entity := &TestEntity{}
	entity.PartitionKey = "ABC"
	entity.RowKey = "123"
	entity.Property1 = "Value1"
	entity.Property2 = "Value2"
	entity.Property3 = "Value3"

	json, _ := json.Marshal(entity)
	tableStorageProxy.MergeEntity(Table, "ABC", "456", json)
}

func TestInsertOrMergeEntity(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	entity := &TestEntity{}
	entity.PartitionKey = "ABC"
	entity.RowKey = "123"
	entity.Property1 = "Value1"
	entity.Property2 = "Value2"
	entity.Property3 = "Value3"

	json, _ := json.Marshal(entity)
	tableStorageProxy.InsertOrMergeEntity(Table, "ABC", "456", json)
}

func TestInsertOrReplaceEntity(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	entity := &TestEntity{}
	entity.PartitionKey = "ABC"
	entity.RowKey = "123"
	entity.Property1 = "Value1"
	entity.Property2 = "Value2"
	entity.Property3 = "Value3"

	json, _ := json.Marshal(entity)
	tableStorageProxy.InsertOrReplaceEntity(Table, "ABC", "456", json)
}

func TestDeleteTable(t *testing.T) {
	goHaveStorage := New(Account, Key)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()
	tableStorageProxy.DeleteTable(Table)
}

type TestEntity struct {
	PartitionKey string
	RowKey       string
	Property1    string
	Property2    string
	Property3    string
}
