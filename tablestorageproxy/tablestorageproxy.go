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
	baseUrl string
}

func New(goHaveStorage GoHaveStorage) *TableStorageProxy {
	var tableStorageProxy TableStorageProxy

	tableStorageProxy.goHaveStorage = goHaveStorage
	tableStorageProxy.baseUrl = "https://"+goHaveStorage.GetAccount()+".table.core.windows.net/"

	return &tableStorageProxy
}

func (tableStorageProxy *TableStorageProxy) QueryTables() {
	tableStorageProxy.get("Tables", "")
}

func (tableStorageProxy *TableStorageProxy) QueryEntity(tableName string, partitionKey string, rowKey string, selects string) {
	tableStorageProxy.get(tableName + "%28PartitionKey=%27" + partitionKey + "%27,RowKey=%27" + rowKey + "%27%29", "?$select="+selects)
}

func (tableStorageProxy *TableStorageProxy) QueryEntities(tableName string, selects string, filter string, top string) {
	tableStorageProxy.get(tableName, "?$filter="+filter + "&$select=" + selects+"&$top="+top)
}

func (tableStorageProxy *TableStorageProxy) DeleteEntity(tableName string, partitionKey string, rowKey string) {
	tableStorageProxy.executeEntityRequest("DELETE",tableName, partitionKey, rowKey, nil, true)
}

func (tableStorageProxy *TableStorageProxy) UpdateEntity(tableName string, partitionKey string, rowKey string, json []byte) {
	tableStorageProxy.executeEntityRequest("PUT",tableName, partitionKey, rowKey, json, true)
}

func (tableStorageProxy *TableStorageProxy) MergeEntity(tableName string, partitionKey string, rowKey string, json []byte) {
	tableStorageProxy.executeEntityRequest("MERGE",tableName, partitionKey, rowKey, json, true)
}

func (tableStorageProxy *TableStorageProxy) InsertOrMergeEntity(tableName string, partitionKey string, rowKey string, json []byte) {
	tableStorageProxy.executeEntityRequest("MERGE",tableName, partitionKey, rowKey, json, false)
}

func (tableStorageProxy *TableStorageProxy) InsertOrReplaceEntity(tableName string, partitionKey string, rowKey string, json []byte) {
	tableStorageProxy.executeEntityRequest("PUT",tableName, partitionKey, rowKey, json, false)
}

func (tableStorageProxy *TableStorageProxy) DeleteTable(tableName string) {
	target := "Tables%28%27" + tableName + "%27%29"

	client := &http.Client{}
	request, _ := http.NewRequest("DELETE", tableStorageProxy.baseUrl+target, nil)
	request.Header.Set("Content-Type", "application/atom+xml")

	tableStorageProxy.executeRequest(request, client, target)
}

type CreateTableArgs struct {
	TableName string
}

func (tableStorageProxy *TableStorageProxy) CreateTable(tableName string) {
	var createTableArgs CreateTableArgs
	createTableArgs.TableName = tableName

	json, _ := json.Marshal(createTableArgs)
	tableStorageProxy.postJson("Tables", json)
}

func (tableStorageProxy *TableStorageProxy) InsertEntity(tableName string, json []byte) {
	tableStorageProxy.postJson(tableName, json)
}

func (tableStorageProxy *TableStorageProxy) get(target string, query string) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", tableStorageProxy.baseUrl+target + query, nil)
	request.Header.Set("Accept", "application/json;odata=nometadata")

	tableStorageProxy.executeRequest(request, client, target)
}

func (tableStorageProxy *TableStorageProxy) postJson(target string, json []byte) {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", tableStorageProxy.baseUrl+target, bytes.NewBuffer(json))
	addPayloadHeaders(request, len(json))

	tableStorageProxy.executeRequest(request, client, target)
}

func addPayloadHeaders(request *http.Request, bodyLength int) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Length", string(bodyLength))
}

func (tableStorageProxy *TableStorageProxy) executeRequest(request *http.Request, client *http.Client, target string) {
	xmsdate, Authentication := tableStorageProxy.calculateDateAndAuthentication(target)

	request.Header.Set("x-ms-date", xmsdate)
	request.Header.Set("x-ms-version", "2013-08-15")
	request.Header.Set("Authorization", Authentication)

	requestDump, _ := httputil.DumpRequest(request, true)

	fmt.Printf("Request: %s\n", requestDump)

	response, _ := client.Do(request)

	responseDump, _ := httputil.DumpResponse(response, true)
	fmt.Printf("Response: %s\n", responseDump)
}

func (tableStorageProxy *TableStorageProxy)  executeEntityRequest(httpVerb string, tableName string, partitionKey string, rowKey string, json []byte, useIfMatch bool) {
	client := &http.Client{}
	request, _ := http.NewRequest(httpVerb, tableStorageProxy.baseUrl+tableName + "%28PartitionKey=%27" + partitionKey + "%27,RowKey=%27" + rowKey + "%27%29",  bytes.NewBuffer(json))

	if json != nil {
		addPayloadHeaders(request, len(json))
	}

	if useIfMatch {
		request.Header.Set("If-Match", "*")
	}

	tableStorageProxy.executeRequest(request, client, tableName + "%28PartitionKey=%27" + partitionKey + "%27,RowKey=%27" + rowKey + "%27%29")
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
