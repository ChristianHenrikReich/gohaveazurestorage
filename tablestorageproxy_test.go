package gohavestorage

import (
	"encoding/json"
	"fmt"
	"gohavestorage/gohavestoragecommon"
	"reflect"
	"strings"
	"testing"
)

var Key = ""
var Account = ""
var Table = "TestTable"

func TestTableMethods(t *testing.T) {
	table := "TableForTestingTableMethods"

	goHaveStorage := NewWithDebug(Account, Key, false)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	httpStatusCode := tableStorageProxy.CreateTable(table)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	body, httpStatusCode := tableStorageProxy.QueryTables()
	assertHTTPStatusCode(t, httpStatusCode, 200)
	if strings.Contains(string(body), "\"TableName\":\"TableForTestingTableMethods\"") != true {
		t.Fail()
	}

	httpStatusCode = tableStorageProxy.DeleteTable(table)
	assertHTTPStatusCode(t, httpStatusCode, 204)
}

func TestEntityMethods(t *testing.T) {
	table := "TableForTestingEntityMethods"

	goHaveStorage := NewWithDebug(Account, Key, false)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	httpStatusCode := tableStorageProxy.CreateTable(table)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	jsonBytes, _ := json.Marshal(&TestEntity{"ABC", "123", "Value1", "Value2", "Value3"})
	httpStatusCode = tableStorageProxy.InsertEntity(table, jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "456", "Value1", "Value2", "Value3"})
	httpStatusCode = tableStorageProxy.InsertEntity(table, jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "789", "Value1", "Value2", "Value3"})
	httpStatusCode = tableStorageProxy.InsertEntity(table, jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	body, httpStatusCode := tableStorageProxy.QueryEntities(table, "PartitionKey,RowKey,Property1,Property2,Property3", "RowKey gt '123'", "1")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"value\":[{\"PartitionKey\":\"ABC\",\"RowKey\":\"456\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value3\"}]}")

	httpStatusCode = tableStorageProxy.DeleteEntity(table, "ABC", "456")
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorageProxy.QueryEntities(table, "PartitionKey,RowKey,Property1,Property2,Property3", "", "")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"value\":[{\"PartitionKey\":\"ABC\",\"RowKey\":\"123\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value3\"},{\"PartitionKey\":\"ABC\",\"RowKey\":\"789\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value3\"}]}")

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "456", "Value1", "Value2", "Value333"})
	httpStatusCode = tableStorageProxy.InsertOrReplaceEntity(table, "ABC", "456", jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorageProxy.QueryEntity(table, "ABC", "456", "PartitionKey,RowKey,Property1,Property2,Property3")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"PartitionKey\":\"ABC\",\"RowKey\":\"456\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value333\"}")

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "456", "Value1", "Value2", "Value3"})
	httpStatusCode = tableStorageProxy.InsertOrReplaceEntity(table, "ABC", "456", jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorageProxy.QueryEntity(table, "ABC", "456", "PartitionKey,RowKey,Property1,Property2,Property3")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"PartitionKey\":\"ABC\",\"RowKey\":\"456\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value3\"}")

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "456", "Value1", "Value2", "Value333"})
	httpStatusCode = tableStorageProxy.UpdateEntity(table, "ABC", "456", jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorageProxy.QueryEntity(table, "ABC", "456", "PartitionKey,RowKey,Property1,Property2,Property3")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"PartitionKey\":\"ABC\",\"RowKey\":\"456\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value333\"}")

	jsonBytes, _ = json.Marshal(&SmallTestEntity{PartitionKey: "ABC", RowKey: "246", Property1: "Value1"})
	httpStatusCode = tableStorageProxy.InsertOrMergeEntity(table, "ABC", "246", jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorageProxy.QueryEntity(table, "ABC", "246", "PartitionKey,RowKey,Property1,Property2,Property3")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"PartitionKey\":\"ABC\",\"RowKey\":\"246\",\"Property1\":\"Value1\",\"Property2\":null,\"Property3\":null}")

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "246", "Value1", "Value2", "Value3"})
	httpStatusCode = tableStorageProxy.MergeEntity(table, "ABC", "246", jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorageProxy.QueryEntity(table, "ABC", "246", "PartitionKey,RowKey,Property1,Property2,Property3")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"PartitionKey\":\"ABC\",\"RowKey\":\"246\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value3\"}")

	httpStatusCode = tableStorageProxy.DeleteTable(table)
	assertHTTPStatusCode(t, httpStatusCode, 204)
}

func assertHTTPStatusCode(t *testing.T, httpStatusCode int, expected int) {
	if httpStatusCode != expected {
		fmt.Printf("Faild http code other than expected:%d\n", httpStatusCode)
		t.Fail()
	}
}

func assertBody(t *testing.T, body []byte, expected string) {
	if string(body) != expected {
		fmt.Printf("Unexpected return:\n%s\nvs\n%s\n", string(body), expected)
		t.Fail()
	}
}

func TestTableServiceProperties(t *testing.T) {
	goHaveStorage := NewWithDebug(Account, Key, false)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()
	properties, _ := tableStorageProxy.GetTableServiceProperties()
	httpStatusCode := tableStorageProxy.SetTableServiceProperties(properties)

	lastestProperties, _ := tableStorageProxy.GetTableServiceProperties()

	if httpStatusCode != 202 {
		fmt.Printf("Faild http code other than expected:%d", httpStatusCode)
		t.Fail()
	}
	if reflect.DeepEqual(properties, lastestProperties) == false {
		fmt.Printf("Dump:\n%+v\n\nvs\n\n%+v", properties, lastestProperties)
		t.Fail()
	}
}

func TestGetTableServiceStats(t *testing.T) {
	goHaveStorage := NewWithDebug(Account, Key, false)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()
	stats, httpStatusCode := tableStorageProxy.GetTableServiceStats()

	if httpStatusCode != 200 {
		fmt.Printf("Faild http code other than expected:%d", httpStatusCode)
		t.Fail()
	}
	if stats.GeoReplication.Status == "" || stats.GeoReplication.LastSyncTime == "" {
		t.Fail()
	}
}

func TestTableACL(t *testing.T) {
	goHaveStorage := NewWithDebug(Account, Key, false)
	tableStorageProxy := goHaveStorage.NewTableStorageProxy()

	accessPolicy := gohavestoragecommon.AccessPolicy{Start: "2014-12-31T00:00:00.0000000Z", Expiry: "2114-12-31T00:00:00.0000000Z", Permission: "raud"}
	signedIdentifier := gohavestoragecommon.SignedIdentifier{Id: "b54df8ab0e2d52759110f48c8d0c19e2", AccessPolicy: accessPolicy}
	signedIdentifiers := &gohavestoragecommon.SignedIdentifiers{[]gohavestoragecommon.SignedIdentifier{signedIdentifier}}
	tableStorageProxy.SetTableACL(Table, signedIdentifiers)

	acl, httpStatusCode := tableStorageProxy.GetTableACL(Table)

	if httpStatusCode != 200 {
		fmt.Printf("Faild http code other than expected:%d", httpStatusCode)
		t.Fail()
	}
	if reflect.DeepEqual(signedIdentifiers, acl) == false {
		fmt.Printf("Dump:\n%+v\n\nvs\n\n%+v", signedIdentifiers, acl)
		t.Fail()
	}
}

type TestEntity struct {
	PartitionKey string
	RowKey       string
	Property1    string
	Property2    string
	Property3    string
}

type SmallTestEntity struct {
	PartitionKey string
	RowKey       string
	Property1    string
}
