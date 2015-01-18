package tablestorageproxy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
	"strings"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"bytes"
	)

type GoHaveStorage interface {
	GetKey()	[]byte
	GetAccount() string
}

type TableStorageProxy struct {
	goHaveStorage GoHaveStorage
}

func New(goHaveStorage GoHaveStorage) *TableStorageProxy {
	var tableStorageProxy TableStorageProxy

	tableStorageProxy.goHaveStorage = goHaveStorage

	return &tableStorageProxy
}

func (tableStorageProxy *TableStorageProxy) QueryTables() {
	xmsdate, Authentication := tableStorageProxy.calculateDateAndAuthentication("Tables")

	client := &http.Client{}
		request, _ := http.NewRequest("GET", "https://" + tableStorageProxy.goHaveStorage.GetAccount() + ".table.core.windows.net/Tables", nil)
		request.Header.Set("Accept","application/json;odata=nometadata")

		tableStorageProxy.executeRequest(request, client)
	}

func (tableStorageProxy *TableStorageProxy) DeleteTable(tableName string) {
	xmsdate, Authentication := tableStorageProxy.calculateDateAndAuthentication("Tables%28%27" + tableName + "%27%29")

	client := &http.Client{}
	request, _ := http.NewRequest("DELETE", "https://" + tableStorageProxy.goHaveStorage.GetAccount() + ".table.core.windows.net/Tables('" + tableName + "')", nil)
	request.Header.Set("Content-Type","application/atom+xml")

	tableStorageProxy.executeRequest(request, client)
}

type CreateTableArgs struct {
	TableName string
}

func (tableStorageProxy *TableStorageProxy) CreateTable(tableName string) {
	var createTableArgs CreateTableArgs
	createTableArgs.TableName = tableName

	jsonBytes, _ := json.Marshal(createTableArgs)

	xmsdate, Authentication := tableStorageProxy.calculateDateAndAuthentication("Tables")

	client := &http.Client{}
	request, _ := http.NewRequest("POST", "https://" + tableStorageProxy.goHaveStorage.GetAccount() + ".table.core.windows.net/Tables", bytes.NewBuffer(jsonBytes))
	request.Header.Set("Accept","application/json;odata=nometadata")
	request.Header.Set("Content-Type","application/json")
	request.Header.Set("Content-Length",string(len(jsonBytes)))

	tableStorageProxy.executeRequest(request, client)
}

func (tableStorageProxy *TableStorageProxy) executeRequest(request *http.Request, client *http.Client) {
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
