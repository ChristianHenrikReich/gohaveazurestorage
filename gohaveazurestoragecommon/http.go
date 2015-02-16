package gohaveazurestoragecommon

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"
	"time"
)

type HTTP struct {
	baseURL          string
	secondaryBaseURL string
	account          string
	key              []byte
	dumpSessions     bool
}

func NewHTTP(storageType string, account string, key []byte, dumpSessions bool) *HTTP {
	http := &HTTP{account: account, key: key, dumpSessions: dumpSessions}
	http.baseURL = "https://" + account + "." + storageType + ".core.windows.net/"
	http.secondaryBaseURL = "https://" + account + "-secondary." + storageType + ".core.windows.net/"

	return http
}

func (storagehttp *HTTP) Request(httpVerb string, target string, query string, json []byte, headers map[string]string, useAccept bool, useContentTypeXML bool, useSecondary bool) ([]byte, int) {
	storagehttp.calculateDateAndAuthentication(httpVerb, target, headers)

	baseURL := ""
	if useSecondary {
		baseURL = storagehttp.secondaryBaseURL
	} else {
		baseURL = storagehttp.baseURL
	}

	client := &http.Client{}
	request, _ := http.NewRequest(httpVerb, baseURL+target+query, bytes.NewBuffer(json))

	if json != nil {
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Length", string(len(json)))
	}

	if useContentTypeXML {
		request.Header.Set("Content-Type", "application/xml")
	}

	if useAccept {
		request.Header.Set("Accept", "application/json;odata=nometadata")
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if storagehttp.dumpSessions {
		responseDump, _ := httputil.DumpResponse(response, true)
		requestDump, _ := httputil.DumpRequest(request, true)

		fmt.Printf("Request: %s\n", requestDump)
		fmt.Printf("%s\n", string(json))
		fmt.Printf("Response: %s\n", responseDump)
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	return contents, response.StatusCode
}

func (storagehttp *HTTP) calculateDateAndAuthentication(httpVerb string, target string, headers map[string]string) {
	xmsdate := strings.Replace(time.Now().UTC().Add(-time.Minute).Format(time.RFC1123), "UTC", "GMT", -1)

	headers["x-ms-date"] = xmsdate
	headers["x-ms-version"] = "2013-08-15"

	signatureString := xmsdate + "\n/" + storagehttp.account + "/" + target
	headers["Authorization"] = "SharedKeyLite " + storagehttp.account + ":" + computeHmac256(signatureString, storagehttp.key)
}

func (storagehttp *HTTP) calculateDateAndAuthenticationForBlobs(httpVerb string, target string, headers map[string]string) {
	headers["x-ms-date"] = strings.Replace(time.Now().UTC().Add(-time.Minute).Format(time.RFC1123), "UTC", "GMT", -1)
	headers["x-ms-version"] = "2013-08-15"

	var keys []string
	for key := range headers {
		if strings.HasPrefix(key, "x-ms-") {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	conicalHeaders := ""
	for _, key := range keys {
		conicalHeaders += key + ":" + headers[key] + "\n"
	}

	signatureString := httpVerb + "\n\n\n\n" + conicalHeaders + "/" + storagehttp.account + "/" + target
	headers["Authorization"] = "SharedKeyLite " + storagehttp.account + ":" + computeHmac256(signatureString, storagehttp.key)
}

func computeHmac256(message string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
