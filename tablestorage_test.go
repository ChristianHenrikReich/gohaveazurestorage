package gohaveazurestorage

import (
	"encoding/json"
	"fmt"
	"gohaveazurestorage/gohaveazurestoragecommon"
	"reflect"
	"strings"
	"testing"
)

var Key = ""
var Account = ""

func TestTableMethods(t *testing.T) {
	table := "TableForTestingTableMethods"

	goHaveStorage := NewWithDebug(Account, Key, false)
	tableStorage := goHaveStorage.NewTableStorage()

	httpStatusCode := tableStorage.CreateTable(table)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	body, httpStatusCode := tableStorage.QueryTables()
	assertHTTPStatusCode(t, httpStatusCode, 200)
	if strings.Contains(string(body), "\"TableName\":\"TableForTestingTableMethods\"") != true {
		t.Fail()
	}

	httpStatusCode = tableStorage.DeleteTable(table)
	assertHTTPStatusCode(t, httpStatusCode, 204)
}

func TestEntityMethods(t *testing.T) {
	table := "TableForTestingEntityMethods"

	goHaveStorage := NewWithDebug(Account, Key, false)
	tableStorage := goHaveStorage.NewTableStorage()

	httpStatusCode := tableStorage.CreateTable(table)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	jsonBytes, _ := json.Marshal(&TestEntity{"ABC", "123", "Value1", "Value2", "Value3"})
	httpStatusCode = tableStorage.InsertEntity(table, jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "456", "Value1", "Value2", "Value3"})
	httpStatusCode = tableStorage.InsertEntity(table, jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "789", "Value1", "Value2", "Value3"})
	httpStatusCode = tableStorage.InsertEntity(table, jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	body, httpStatusCode := tableStorage.QueryEntities(table, "PartitionKey,RowKey,Property1,Property2,Property3", "RowKey gt '123'", "1")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"value\":[{\"PartitionKey\":\"ABC\",\"RowKey\":\"456\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value3\"}]}")

	httpStatusCode = tableStorage.DeleteEntity(table, "ABC", "456")
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorage.QueryEntities(table, "PartitionKey,RowKey,Property1,Property2,Property3", "", "")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"value\":[{\"PartitionKey\":\"ABC\",\"RowKey\":\"123\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value3\"},{\"PartitionKey\":\"ABC\",\"RowKey\":\"789\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value3\"}]}")

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "456", "Value1", "Value2", "Value333"})
	httpStatusCode = tableStorage.InsertOrReplaceEntity(table, "ABC", "456", jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorage.QueryEntity(table, "ABC", "456", "PartitionKey,RowKey,Property1,Property2,Property3")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"PartitionKey\":\"ABC\",\"RowKey\":\"456\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value333\"}")

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "456", "Value1", "Value2", "Value3"})
	httpStatusCode = tableStorage.InsertOrReplaceEntity(table, "ABC", "456", jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorage.QueryEntity(table, "ABC", "456", "PartitionKey,RowKey,Property1,Property2,Property3")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"PartitionKey\":\"ABC\",\"RowKey\":\"456\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value3\"}")

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "456", "Value1", "Value2", "Value333"})
	httpStatusCode = tableStorage.UpdateEntity(table, "ABC", "456", jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorage.QueryEntity(table, "ABC", "456", "PartitionKey,RowKey,Property1,Property2,Property3")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"PartitionKey\":\"ABC\",\"RowKey\":\"456\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value333\"}")

	jsonBytes, _ = json.Marshal(&SmallTestEntity{PartitionKey: "ABC", RowKey: "246", Property1: "Value1"})
	httpStatusCode = tableStorage.InsertOrMergeEntity(table, "ABC", "246", jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorage.QueryEntity(table, "ABC", "246", "PartitionKey,RowKey,Property1,Property2,Property3")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"PartitionKey\":\"ABC\",\"RowKey\":\"246\",\"Property1\":\"Value1\",\"Property2\":null,\"Property3\":null}")

	jsonBytes, _ = json.Marshal(&TestEntity{"ABC", "246", "Value1", "Value2", "Value3"})
	httpStatusCode = tableStorage.MergeEntity(table, "ABC", "246", jsonBytes)
	assertHTTPStatusCode(t, httpStatusCode, 204)

	body, httpStatusCode = tableStorage.QueryEntity(table, "ABC", "246", "PartitionKey,RowKey,Property1,Property2,Property3")
	assertHTTPStatusCode(t, httpStatusCode, 200)
	assertBody(t, body, "{\"PartitionKey\":\"ABC\",\"RowKey\":\"246\",\"Property1\":\"Value1\",\"Property2\":\"Value2\",\"Property3\":\"Value3\"}")

	httpStatusCode = tableStorage.DeleteTable(table)
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
	tableStorage := goHaveStorage.NewTableStorage()
	properties, _ := tableStorage.GetTableServiceProperties()
	httpStatusCode := tableStorage.SetTableServiceProperties(properties)

	lastestProperties, _ := tableStorage.GetTableServiceProperties()

	assertHTTPStatusCode(t, httpStatusCode, 202)
	if reflect.DeepEqual(properties, lastestProperties) == false {
		fmt.Printf("Dump:\n%+v\n\nvs\n\n%+v", properties, lastestProperties)
		t.Fail()
	}
}

func TestGetTableServiceStats(t *testing.T) {
	goHaveStorage := NewWithDebug(Account, Key, false)
	tableStorage := goHaveStorage.NewTableStorage()
	stats, httpStatusCode := tableStorage.GetTableServiceStats()

	assertHTTPStatusCode(t, httpStatusCode, 200)
	if stats.GeoReplication.Status == "" || stats.GeoReplication.LastSyncTime == "" {
		t.Fail()
	}
}

func TestTableACL(t *testing.T) {
	table := "TableForTestingACLMethods"

	goHaveStorage := NewWithDebug(Account, Key, false)
	tableStorage := goHaveStorage.NewTableStorage()

	httpStatusCode := tableStorage.CreateTable(table)
	assertHTTPStatusCode(t, httpStatusCode, 201)

	accessPolicy := gohaveazurestoragecommon.AccessPolicy{Start: "2014-12-31T00:00:00.0000000Z", Expiry: "2114-12-31T00:00:00.0000000Z", Permission: "raud"}
	signedIdentifier := gohaveazurestoragecommon.SignedIdentifier{Id: "b54df8ab0e2d52759110f48c8d0c19e2", AccessPolicy: accessPolicy}
	signedIdentifiers := &gohaveazurestoragecommon.SignedIdentifiers{[]gohaveazurestoragecommon.SignedIdentifier{signedIdentifier}}
	tableStorage.SetTableACL(table, signedIdentifiers)

	acl, httpStatusCode := tableStorage.GetTableACL(table)

	assertHTTPStatusCode(t, httpStatusCode, 200)
	if reflect.DeepEqual(signedIdentifiers, acl) == false {
		fmt.Printf("Dump:\n%+v\n\nvs\n\n%+v", signedIdentifiers, acl)
		t.Fail()
	}

	httpStatusCode = tableStorage.DeleteTable(table)
	assertHTTPStatusCode(t, httpStatusCode, 204)
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
