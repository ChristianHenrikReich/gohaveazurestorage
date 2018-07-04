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

func (storagehttp *HTTP) Request(httpVerb string, target string, query string, json []byte, useIfMatch bool, useAccept bool, useContentTypeXML bool, useSecondary bool) ([]byte, int) {
	xmsdate, Authentication := storagehttp.calculateDateAndAuthentication(target)

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

	if useIfMatch {
		request.Header.Set("If-Match", "*")
	}

	if useAccept {
		request.Header.Set("Accept", "application/json;odata=nometadata")
	}

	request.Header.Set("x-ms-date", xmsdate)
	request.Header.Set("x-ms-version", "2013-08-15")
	request.Header.Set("Authorization", Authentication)

	response, err := client.Do(request)
	if err != nil {
		return nil, http.StatusServiceUnavailable
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
		return nil, http.StatusUnprocessableEntity
	}

	return contents, response.StatusCode
}

func (storagehttp *HTTP) calculateDateAndAuthentication(target string) (string, string) {
	xmsdate := strings.Replace(time.Now().UTC().Add(-time.Minute).Format(time.RFC1123), "UTC", "GMT", -1)
	SignatureString := xmsdate + "\n/" + storagehttp.account + "/" + target
	Authentication := "SharedKeyLite " + storagehttp.account + ":" + computeHmac256(SignatureString, storagehttp.key)
	return xmsdate, Authentication
}

func computeHmac256(message string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
