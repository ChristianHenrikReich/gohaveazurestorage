package tablestorageproxy

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type GoHaveStorage interface {
	GetKey() []byte
	GetAccount() string
}

type TableStorageProxy struct {
	goHaveStorage GoHaveStorage
	baseUrl       string
}

func New(goHaveStorage GoHaveStorage) *TableStorageProxy {
	var tableStorageProxy TableStorageProxy

	tableStorageProxy.goHaveStorage = goHaveStorage
	tableStorageProxy.baseUrl = "https://" + goHaveStorage.GetAccount() + ".table.core.windows.net/"

	return &tableStorageProxy
}

func (tableStorageProxy *TableStorageProxy) QueryTables() {
	tableStorageProxy.executeCommonRequest("GET", "Tables", "", nil, false, true, false)
}

func (tableStorageProxy *TableStorageProxy) QueryEntity(tableName string, partitionKey string, rowKey string, selects string) {
	tableStorageProxy.executeCommonRequest("GET", tableName+"%28PartitionKey=%27"+partitionKey+"%27,RowKey=%27"+rowKey+"%27%29", "?$select="+selects, nil, false, true, false)
}

func (tableStorageProxy *TableStorageProxy) QueryEntities(tableName string, selects string, filter string, top string) {
	tableStorageProxy.executeCommonRequest("GET", tableName, "?$filter="+filter+"&$select="+selects+"&$top="+top, nil, false, true, false)
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

func (tableStorageProxy *TableStorageProxy) DeleteTable(tableName string) {
	tableStorageProxy.executeCommonRequest("DELETE", tableName, "", nil, false, false, true)
}

type CreateTableArgs struct {
	TableName string
}

func (tableStorageProxy *TableStorageProxy) CreateTable(tableName string) {
	json, _ := json.Marshal(CreateTableArgs{TableName: tableName})
	tableStorageProxy.executeCommonRequest("POST", "Tables", "", json, false, true, false)
}

func (tableStorageProxy *TableStorageProxy) InsertEntity(tableName string, json []byte) {
	tableStorageProxy.executeCommonRequest("POST", tableName, "", json, false, true, false)
}

func (tableStorageProxy *TableStorageProxy) executeEntityRequest(httpVerb string, tableName string, partitionKey string, rowKey string, json []byte, useIfMatch bool) {
	tableStorageProxy.executeCommonRequest(httpVerb, tableName+"%28PartitionKey=%27"+partitionKey+"%27,RowKey=%27"+rowKey+"%27%29", "", json, useIfMatch, false, false)
}

func (tableStorageProxy *TableStorageProxy) executeCommonRequest(httpVerb string, target string, query string, json []byte, useIfMatch bool, useAccept bool, useContentTypeXML bool) {
	xmsdate, Authentication := tableStorageProxy.calculateDateAndAuthentication(target)

	client := &http.Client{}
	request, _ := http.NewRequest(httpVerb, tableStorageProxy.baseUrl+target+query, bytes.NewBuffer(json))

	if json != nil {
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Length", string(len(json)))
	}

	if useContentTypeXML {
		request.Header.Set("Content-Type", "application/atom+xml")
	}

	if useIfMatch {
		request.Header.Set("If-Match", "*")
	}

	if useAccept {
		request.Header.Set("Accept", "application/json;odata=nometadata")
	}

	request.Header.Set("x-ms-date", xmsdate)
	request.Header.Set("x-ms-version", "2013-08-15")
	request.Header.Set("Authorization", Authentication)

	requestDump, _ := httputil.DumpRequest(request, true)

	fmt.Printf("Request: %s\n", requestDump)

	response, _ := client.Do(request)

	responseDump, _ := httputil.DumpResponse(response, true)
	fmt.Printf("Response: %s\n", responseDump)
}

func (tableStorageProxy *TableStorageProxy) calculateDateAndAuthentication(target string) (string, string) {
	xmsdate := strings.Replace(time.Now().UTC().Add(-time.Minute).Format(time.RFC1123), "UTC", "GMT", -1)
	SignatureString := xmsdate + "\n/" + tableStorageProxy.goHaveStorage.GetAccount() + "/" + target
	Authentication := "SharedKeyLite " + tableStorageProxy.goHaveStorage.GetAccount() + ":" + computeHmac256(SignatureString, tableStorageProxy.goHaveStorage.GetKey())
	return xmsdate, Authentication
}

func computeHmac256(message string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
