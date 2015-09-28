// Copyright (c) 2015 RightScale Inc, All Rights Reserved

package demo

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rightscale/go-boilerplate/misc"
	"github.com/rightscale/gojiutil"
	"github.com/zenazn/goji/web"
	"gopkg.in/inconshreveable/log15.v2"
)

// This demo file shows two ways to test the handlers, it is not suggested to use both in a real
// project, the two methods are provided here as a sample.

// Tests that simply create a mux and exercise that by calling the Mux's ServeHTTP function
// This is great for isolated tests, it becomes difficult when the middleware higher in the stack
// is needed or other concurrent goroutines and other handlers al also needed for higher-level tests
var _ = Describe("Mux-based settings tests", func() {

	var mx *web.Mux // mux with the handlers we're testing

	BeforeEach(func() {
		settings = make(map[string]string)
		mx = NewMux()
		gojiutil.AddCommon15(mx, log15.Root())
		mx.Use(gojiutil.ParamsLogger(true)) // useful for troubleshooting
	})

	It("gets what it sets", func() {
		// set a value
		req, _ := http.NewRequest("PUT", "http://example.com/settings/hello?value=world",
			bytes.NewReader([]byte{}))
		resp := httptest.NewRecorder()
		mx.ServeHTTP(resp, req)
		Ω(resp.Code).Should(Equal(200))
		Ω(settings["hello"]).Should(Equal("world"))

		// get the value back
		req, _ = http.NewRequest("GET", "http://example.com/settings/hello", nil)
		resp = httptest.NewRecorder()
		mx.ServeHTTP(resp, req)
		Ω(resp.Code).Should(Equal(200))
		Ω(resp.Body.String()).Should(Equal("world"))
	})

	It("deletes and lists", func() {
		// set a value
		req, _ := http.NewRequest("PUT", "http://example.com/settings/hello?value=world",
			bytes.NewReader([]byte{}))
		resp := httptest.NewRecorder()
		mx.ServeHTTP(resp, req)
		Ω(resp.Code).Should(Equal(200))
		Ω(settings["hello"]).Should(Equal("world"))
		Ω(settings).Should(HaveLen(1))

		// list the values
		req, _ = http.NewRequest("GET", "http://example.com/settings", nil)
		resp = httptest.NewRecorder()
		mx.ServeHTTP(resp, req)
		Ω(resp.Code).Should(Equal(200))
		Ω(resp.Body.String()).Should(Equal(`{"hello":"world"}`))

		// delete the value
		req, _ = http.NewRequest("DELETE", "http://example.com/settings/hello", nil)
		resp = httptest.NewRecorder()
		mx.ServeHTTP(resp, req)
		Ω(resp.Code).Should(Equal(201))
		Ω(settings).Should(HaveLen(0))
	})
})

// Tests that create a web server and run http requests against it.
// This is great because one can make a request with a 1-liner and all the middleware is
// right there, so it's a true representation of what a client will do. Also it's easy to
// invoke handlers outside of this package and to have background goroutines running (they
// just need to have a way to be terminated in the AfterEach block)
var _ = Describe("HTTP-based settings tests", func() {

	var server *httptest.Server

	BeforeEach(func() {
		settings = make(map[string]string)
		mx := NewMux()
		gojiutil.AddCommon15(mx, log15.Root())
		mx.Use(gojiutil.ParamsLogger(true)) // useful for troubleshooting
		server = httptest.NewServer(mx)
	})

	AfterEach(func() {
		server.Close()
	})

	It("gets what it sets", func() {
		// MakeRequest issues a request and checks the response status code
		respBody, _ := misc.MakeRequest("PUT", server.URL+"/settings/hello?value=world",
			"", 200)
		Ω(settings["hello"]).Should(Equal("world"))
		respBody, _ = misc.MakeRequest("GET", server.URL+"/settings/hello", "", 200)
		Ω(respBody).Should(Equal("world"))
	})

	It("deletes and lists", func() {
		// MakeRequest issues a request and checks the response status code
		misc.MakeRequest("PUT", server.URL+"/settings/hello?value=world", "", 200)
		Ω(settings["hello"]).Should(Equal("world"))

		respObj, _ := misc.MakeRequestObj("GET", server.URL+"/settings", "", 200)
		Ω(respObj).Should(HaveLen(1))
		Ω(respObj).Should(HaveKeyWithValue("hello", "world"))
	})
})
