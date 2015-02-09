package tablestorage

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gohaveazurestorage/gohaveazurestoragecommon"
	"os"
	"strings"
)

type TableStorage struct {
	http             *gohaveazurestoragecommon.HTTP
	baseUrl          string
	secondaryBaseUrl string
}

func New(http *gohaveazurestoragecommon.HTTP) *TableStorage {
	return &TableStorage{http: http}
}

func (tableStorage *TableStorage) GetTableACL(tableName string) (signedIdentifiers *gohaveazurestoragecommon.SignedIdentifiers, httpStatusCode int) {
	body, httpStatusCode := tableStorage.http.Request("GET", tableName+"?comp=acl", "", nil, false, false, false, false)

	signedIdentifiers = &gohaveazurestoragecommon.SignedIdentifiers{}
	desrializeXML([]byte(body), signedIdentifiers)

	return signedIdentifiers, httpStatusCode
}

func (tableStorage *TableStorage) SetTableACL(tableName string, signedIdentifiers *gohaveazurestoragecommon.SignedIdentifiers) (httpStatusCode int) {
	xmlBytes, _ := xml.MarshalIndent(signedIdentifiers, "", "")
	_, httpStatusCode = tableStorage.http.Request("PUT", tableName+"?comp=acl", "", xmlBytes, false, false, true, false)
	return httpStatusCode
}

func (tableStorage *TableStorage) GetTableServiceProperties() (storageServiceProperties *gohaveazurestoragecommon.StorageServiceProperties, httpStatusCode int) {
	body, httpStatusCode := tableStorage.http.Request("GET", "?comp=properties", "&restype=service", nil, false, false, true, false)

	storageServiceProperties = &gohaveazurestoragecommon.StorageServiceProperties{}
	desrializeXML([]byte(body), storageServiceProperties)
	return storageServiceProperties, httpStatusCode
}

func (tableStorage *TableStorage) SetTableServiceProperties(storageServiceProperties *gohaveazurestoragecommon.StorageServiceProperties) (httpStatusCode int) {
	xmlBytes, _ := xml.MarshalIndent(storageServiceProperties, "", "")
	_, httpStatusCode = tableStorage.http.Request("PUT", "?comp=properties", "&restype=service", append([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"), xmlBytes...), false, false, false, false)
	return httpStatusCode
}

func (tableStorage *TableStorage) GetTableServiceStats() (storageServiceStats *gohaveazurestoragecommon.StorageServiceStats, httpStatusCode int) {
	body, httpStatusCode := tableStorage.http.Request("GET", "?comp=stats", "&restype=service", nil, false, false, false, true)

	storageServiceStats = &gohaveazurestoragecommon.StorageServiceStats{}
	desrializeXML([]byte(body), storageServiceStats)

	return storageServiceStats, httpStatusCode
}

func (tableStorage *TableStorage) QueryEntity(tableName string, partitionKey string, rowKey string, selects string) (body []byte, httpStatusCode int) {
	return tableStorage.http.Request("GET", tableName+"%28PartitionKey=%27"+partitionKey+"%27,RowKey=%27"+rowKey+"%27%29", "?$select="+selects, nil, false, true, false, false)
}

func (tableStorage *TableStorage) QueryEntities(tableName string, selects string, filter string, top string) (body []byte, httpStatusCode int) {
	filter = strings.Replace(filter, " ", "%20", -1)
	return tableStorage.http.Request("GET", tableName, "?$filter="+filter+"&$select="+selects+"&$top="+top, nil, false, true, false, false)
}

func (tableStorage *TableStorage) InsertEntity(tableName string, json []byte) (httpStatusCode int) {
	_, httpStatusCode = tableStorage.http.Request("POST", tableName, "", json, false, true, false, false)
	return httpStatusCode
}

func (tableStorage *TableStorage) DeleteEntity(tableName string, partitionKey string, rowKey string) (httpStatusCode int) {
	_, httpStatusCode = tableStorage.executeEntityRequest("DELETE", tableName, partitionKey, rowKey, nil, true)
	return httpStatusCode
}

func (tableStorage *TableStorage) UpdateEntity(tableName string, partitionKey string, rowKey string, json []byte) (httpStatusCode int) {
	_, httpStatusCode = tableStorage.executeEntityRequest("PUT", tableName, partitionKey, rowKey, json, true)
	return httpStatusCode
}

func (tableStorage *TableStorage) MergeEntity(tableName string, partitionKey string, rowKey string, json []byte) (httpStatusCode int) {
	_, httpStatusCode = tableStorage.executeEntityRequest("MERGE", tableName, partitionKey, rowKey, json, true)
	return httpStatusCode
}

func (tableStorage *TableStorage) InsertOrMergeEntity(tableName string, partitionKey string, rowKey string, json []byte) (httpStatusCode int) {
	_, httpStatusCode = tableStorage.executeEntityRequest("MERGE", tableName, partitionKey, rowKey, json, false)
	return httpStatusCode
}

func (tableStorage *TableStorage) InsertOrReplaceEntity(tableName string, partitionKey string, rowKey string, json []byte) (httpStatusCode int) {
	_, httpStatusCode = tableStorage.executeEntityRequest("PUT", tableName, partitionKey, rowKey, json, false)
	return httpStatusCode
}

type CreateTableArgs struct {
	TableName string
}

func (tableStorage *TableStorage) CreateTable(tableName string) (httpStatusCode int) {
	json, _ := json.Marshal(CreateTableArgs{TableName: tableName})
	_, httpStatusCode = tableStorage.http.Request("POST", "Tables", "", json, false, true, false, false)
	return httpStatusCode
}

func (tableStorage *TableStorage) DeleteTable(tableName string) (httpStatusCode int) {
	_, httpStatusCode = tableStorage.http.Request("DELETE", "Tables%28%27"+tableName+"%27%29", "", nil, false, false, true, false)
	return httpStatusCode
}

func (tableStorage *TableStorage) QueryTables() (body []byte, httpStatusCode int) {
	body, httpStatusCode = tableStorage.http.Request("GET", "Tables", "", nil, false, true, false, false)
	return body, httpStatusCode
}

func (tableStorage *TableStorage) executeEntityRequest(httpVerb string, tableName string, partitionKey string, rowKey string, json []byte, useIfMatch bool) ([]byte, int) {
	return tableStorage.http.Request(httpVerb, tableName+"%28PartitionKey=%27"+partitionKey+"%27,RowKey=%27"+rowKey+"%27%29", "", json, useIfMatch, false, false, false)
}

func desrializeXML(bytes []byte, object interface{}) {
	err := xml.Unmarshal([]byte(bytes), &object)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}
