// Copyright (c) 2015 RightScale Inc, All Rights Reserved

package demo

import (
	"net/http"

	"github.com/rightscale/gojiutil"
	"github.com/zenazn/goji/web"
)

// simple string->string map for demo purposes
var settings = make(map[string]string)

func indexSetting(c web.C, rw http.ResponseWriter, r *http.Request) {
	gojiutil.WriteJSON(c, rw, 200, settings)
}

// getSetting retrieves a setting from the settings map
func getSetting(c web.C, rw http.ResponseWriter, r *http.Request) {
	key := c.URLParams["key"]
	value := settings[key]
	if key == "" || value == "" {
		gojiutil.Errorf(c, rw, 404, `settings key '%s' not found`, key)
		return
	}
	gojiutil.WriteString(rw, 200, value)
}

func putSetting(c web.C, rw http.ResponseWriter, r *http.Request) {
	key := c.URLParams["key"]
	if key == "" {
		gojiutil.ErrorString(c, rw, 413, `settings key missing`)
		return
	}
	value := r.Form.Get("value")
	if key == "" {
		gojiutil.ErrorString(c, rw, 413, `value query string param missing`)
		return
	}
	settings[key] = value
}

func deleteSetting(c web.C, rw http.ResponseWriter, r *http.Request) {
	key := c.URLParams["key"]
	delete(settings, key)
	rw.WriteHeader(201)
}
