// Copyright (c) 2014 RightScale, Inc. - see LICENSE

package misc

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/onsi/gomega"
	"gopkg.in/inconshreveable/log15.v2"
)

// MakeRequest makes a get request, checksthe http status code, and returns the body as string
func MakeRequest(method, url, body string, expectedCode int) (string, *http.Response) {
	log15.Debug("MakeRequest", "verb", method, "url", url)
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, _ := http.NewRequest(method, url, bodyReader)
	resp, err := http.DefaultClient.Do(req)
	gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
	gomega.Ω(resp.StatusCode).Should(gomega.Equal(expectedCode))
	gomega.Ω(resp.Header.Get("Content-Type")).ShouldNot(gomega.BeNil())
	respBody, err := ioutil.ReadAll(resp.Body)
	gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
	return string(respBody), resp
}

// MakeRequestObj makes a request for a JSONobject, checks the http response code, and
// returns the object as map[string]interface{}
func MakeRequestObj(method, url, body string, expectedCode int) (map[string]interface{}, *http.Response) {
	respBody, resp := MakeRequest(method, url, body, expectedCode)
	if respBody == "" {
		return nil, resp
	}
	if resp.StatusCode < 400 {
		gomega.Ω(resp.Header.Get("Content-Type")).Should(gomega.HavePrefix("application/json"))
	} else {
		gomega.Ω(resp.Header.Get("Content-Type")).Should(gomega.HavePrefix("text/plain"))
	}
	// parse json
	var res map[string]interface{}
	err := json.Unmarshal([]byte(respBody), &res)
	gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
	return res, resp
}

// MakeRequestList makes a request for a list of JSON objects, checks the http response code, and
// returns the object list as []map[string]interface{}
func MakeRequestList(method, url string, expectedCode int) ([]map[string]interface{}, *http.Response) {
	respBody, resp := MakeRequest(method, url, "", expectedCode)
	if respBody == "" {
		return nil, resp
	}
	if resp.StatusCode < 400 {
		gomega.Ω(resp.Header.Get("Content-Type")).Should(gomega.HavePrefix("application/json"))
	} else {
		gomega.Ω(resp.Header.Get("Content-Type")).Should(gomega.HavePrefix("text/plain"))
	}
	// parse json
	var res []map[string]interface{}
	err := json.Unmarshal([]byte(respBody), &res)
	gomega.Ω(err).ShouldNot(gomega.HaveOccurred())
	return res, resp
}
