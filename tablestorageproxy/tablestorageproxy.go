package tablestorageproxy

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gohavestorage/gohavestoragecommon"
	"os"
)

type TableStorageProxy struct {
	http             *gohavestoragecommon.HTTP
	baseUrl          string
	secondaryBaseUrl string
}

func New(http *gohavestoragecommon.HTTP) *TableStorageProxy {
	return &TableStorageProxy{http: http}
}

func (tableStorageProxy *TableStorageProxy) GetTableACL(tableName string) (*gohavestoragecommon.SignedIdentifiers, int) {
	body, httpStatusCode := tableStorageProxy.http.Request("GET", tableName+"?comp=acl", "", nil, false, false, false, false)

	response := &gohavestoragecommon.SignedIdentifiers{}
	desrializeXML([]byte(body), response)

	return response, httpStatusCode
}

func (tableStorageProxy *TableStorageProxy) SetTableACL(tableName string, signedIdentifiers *gohavestoragecommon.SignedIdentifiers) int {
	xmlBytes, _ := xml.MarshalIndent(signedIdentifiers, "", "")
	_, httpStatusCode := tableStorageProxy.http.Request("PUT", tableName+"?comp=acl", "", xmlBytes, false, false, true, false)
	return httpStatusCode
}

func (tableStorageProxy *TableStorageProxy) GetTableServiceProperties() (*gohavestoragecommon.StorageServiceProperties, int) {
	body, httpStatusCode := tableStorageProxy.http.Request("GET", "?comp=properties", "&restype=service", nil, false, false, true, false)

	response := &gohavestoragecommon.StorageServiceProperties{}
	desrializeXML([]byte(body), response)
	return response, httpStatusCode
}

func (tableStorageProxy *TableStorageProxy) SetTableServiceProperties(storageServiceProperties *gohavestoragecommon.StorageServiceProperties) int {
	xmlBytes, _ := xml.MarshalIndent(storageServiceProperties, "", "")
	_, httpStatusCode := tableStorageProxy.http.Request("PUT", "?comp=properties", "&restype=service", append([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"), xmlBytes...), false, false, false, false)
	return httpStatusCode
}

func (tableStorageProxy *TableStorageProxy) GetTableServiceStats() (*gohavestoragecommon.StorageServiceStats, int) {
	body, httpStatusCode := tableStorageProxy.http.Request("GET", "?comp=stats", "&restype=service", nil, false, false, false, true)

	response := &gohavestoragecommon.StorageServiceStats{}
	desrializeXML([]byte(body), response)

	return response, httpStatusCode
}

func (tableStorageProxy *TableStorageProxy) QueryEntity(tableName string, partitionKey string, rowKey string, selects string) {
	tableStorageProxy.http.Request("GET", tableName+"%28PartitionKey=%27"+partitionKey+"%27,RowKey=%27"+rowKey+"%27%29", "?$select="+selects, nil, false, true, false, false)
}

func (tableStorageProxy *TableStorageProxy) QueryEntities(tableName string, selects string, filter string, top string) {
	tableStorageProxy.http.Request("GET", tableName, "?$filter="+filter+"&$select="+selects+"&$top="+top, nil, false, true, false, false)
}

func (tableStorageProxy *TableStorageProxy) DeleteEntity(tableName string, partitionKey string, rowKey string) {
	tableStorageProxy.executeEntityRequest("DELETE", tableName, partitionKey, rowKey, nil, true)
}

func (tableStorageProxy *TableStorageProxy) UpdateEntity(tableName string, partitionKey string, rowKey string, json []byte) {
	tableStorageProxy.executeEntityRequest("PUT", tableName, partitionKey, rowKey, json, true)
}

func (tableStorageProxy *TableStorageProxy) MergeEntity(tableName string, partitionKey string, rowKey string, json []byte) {
	tableStorageProxy.executeEntityRequest("MERGE", tableName, partitionKey, rowKey, json, true)
}

func (tableStorageProxy *TableStorageProxy) InsertOrMergeEntity(tableName string, partitionKey string, rowKey string, json []byte) {
	tableStorageProxy.executeEntityRequest("MERGE", tableName, partitionKey, rowKey, json, false)
}

func (tableStorageProxy *TableStorageProxy) InsertOrReplaceEntity(tableName string, partitionKey string, rowKey string, json []byte) {
	tableStorageProxy.executeEntityRequest("PUT", tableName, partitionKey, rowKey, json, false)
}

type CreateTableArgs struct {
	TableName string
}

func (tableStorageProxy *TableStorageProxy) CreateTable(tableName string) int {
	json, _ := json.Marshal(CreateTableArgs{TableName: tableName})
	_, httpStatusCode := tableStorageProxy.http.Request("POST", "Tables", "", json, false, true, false, false)
	return httpStatusCode
}

func (tableStorageProxy *TableStorageProxy) DeleteTable(tableName string) int {
	_, httpStatusCode := tableStorageProxy.http.Request("DELETE", "Tables%28%27"+tableName+"%27%29", "", nil, false, false, true, false)
	return httpStatusCode
}

func (tableStorageProxy *TableStorageProxy) QueryTables() ([]byte, int) {
	body, httpStatusCode := tableStorageProxy.http.Request("GET", "Tables", "", nil, false, true, false, false)
	return body, httpStatusCode
}

func (tableStorageProxy *TableStorageProxy) InsertEntity(tableName string, json []byte) {
	tableStorageProxy.http.Request("POST", tableName, "", json, false, true, false, false)
}

func (tableStorageProxy *TableStorageProxy) executeEntityRequest(httpVerb string, tableName string, partitionKey string, rowKey string, json []byte, useIfMatch bool) {
	tableStorageProxy.http.Request(httpVerb, tableName+"%28PartitionKey=%27"+partitionKey+"%27,RowKey=%27"+rowKey+"%27%29", "", json, useIfMatch, false, false, false)
}

func desrializeXML(bytes []byte, object interface{}) {
	err := xml.Unmarshal([]byte(bytes), &object)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}
