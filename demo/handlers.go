// Copyright (c) 2015 RightScale Inc, All Rights Reserved

package demo

import (
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

// NewMux returns a handler/mux for the demo requests and adds the local handlers to this mux
func NewMux() *web.Mux {
	mx := web.New()
	mx.Use(middleware.SubRouter)
	//mx.Use(getJSONBody)

	// add settings resource
	mx.Get("/settings", indexSetting)
	mx.Get("/settings/:key", getSetting)
	mx.Put("/settings/:key", putSetting)
	mx.Delete("/settings/:key", deleteSetting)

	return mx
}
