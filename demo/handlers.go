// Copyright (c) 2015 RightScale Inc, All Rights Reserved

package demo

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/rightscale/gojiutil"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"gopkg.in/inconshreveable/log15.v2"
)

// NewMux returns a handler/mux for the demo requests and adds the local handlers to this mux
func NewMux() *web.Mux {
	mx := web.New()
	mx.Use(middleware.SubRouter)
	mx.Use(getJSONBody)

	// add settings resource
	mx.Get("/settings/:key", getSetting)
	mx.Put("/settings/:key", putSetting)
	mx.Delete("/settings/:key", deleteSetting)

	return mx
}

// TODO: add to gojiutils

// getJSONBody is a middleware to read and parse an application/json body. This middleware
// is pretty permissive: it allows for having no content-length and no content-type as long
// as either there's no body or the body parses as json it's OK. The result of the json parse is
// stored in c.Env["json"] as a map[string]interface{}, which can be easily mapped to a
// proper struct using github.com/mitchellh/mapstructure
func getJSONBody(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var err error

		// parse content-length header
		cl := 0
		if clh := r.Header.Get("content-length"); clh != "" {
			if cl, err = strconv.Atoi(clh); err != nil {
				gojiutil.ErrorString(*c, rw, 400,
					"Invalid content-length: "+err.Error())
				return
			}
		}

		// parse content-type header
		if ct := r.Header.Get("content-type"); ct != "" && ct != "application/json" {
			gojiutil.ErrorString(*c, rw, 400,
				"Invalid content-type '"+ct+"', application/json expected")
			return
		}

		// try to read body
		var js map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&js)
		switch err {
		case io.EOF:
			if cl != 0 {
				gojiutil.ErrorString(*c, rw, 400, "Premature EOF reading post body")
				return
			}
			log15.Debug("HTTP no request body")
			// got no body, so we're OK
		case nil:
			log15.Info("HTTP Context", "body", js)
			// great!
		default:
			gojiutil.ErrorString(*c, rw, 400, "Cannot parse JSON request body: "+
				err.Error())
			return
		}

		c.Env["json"] = js
		h.ServeHTTP(rw, r)
	})
}
