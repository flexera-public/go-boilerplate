// Copyright (c) 2015 RightScale Inc, All Rights Reserved

package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/rightscale/go-boilerplate/demo"
	"github.com/rightscale/go-boilerplate/log2log15"
	"github.com/rightscale/gojiutil"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/inconshreveable/log15.v2"
)

// FatalError logs and exit app
func FatalError(message string) {
	log15.Crit(message)
	os.Exit(1)
}

// Having a root handler as a global variable is useful if one has to send the std log output
// into log15
var Log15RootHandler = log15.StreamHandler(os.Stdout, log2log15.SimpleFormat(true))

// Embed a version string of the form "go-boilerplate - sdfsfsaf - 24adb4" into the
// executable where..
var VERSION string // Makefile sets this using linker flag, must be uninitialized
func init() {
	// provide a default version string if app is built without makefile
	if VERSION == "" {
		VERSION = "dev-version-manually-built"
	}
}

// Parse command line args and run the thing
func main() {
	// initial logging set-up
	log15.Root().SetHandler(Log15RootHandler)
	log15.Info("RightScale Go Boilerplate", "version", VERSION, "pid", os.Getpid())

	// Parse command line arguments, these flags are often defined as globals so they can be
	// accessed somewhere other than just in the main function
	// You probably want debug to default to true during wild dev and to false once things
	// stabilise, you can use --no-debug to turn it off if the default is true
	debugFlag := kingpin.Flag("debug", "enable debug-level logging").Default("true").Bool()
	kingpin.Version(VERSION)
	kingpin.Parse()

	// further logging set-up, this has some rsc-specific tweaks that may be useful as pattern
	// for other sub-packages
	//rsc_log.Logger = log15.New("pkg", "rsc") // redirect a sub-package's logger
	if !*debugFlag {
		log15.Info("disabling debug logging")
		Log15RootHandler = log15.LvlFilterHandler(log15.LvlInfo, Log15RootHandler)
		log15.Root().SetHandler(Log15RootHandler)
		//rsc_log.Logger.SetHandler(log15.LvlFilterHandler(log15.LvlWarn, Log15RootHandler))
	}
	// more rsc tweaks
	//rsc_httpclient.DumpFormat = rsc_httpclient.NoDump
	//rsc_httpclient.DumpFormat = rsc_httpclient.Debug //+ rsc_httpclient.Verbose
	// redirect the standard log package to log15, some standard library packages log to
	// 'log' directly, for example the http package
	log.SetOutput(log2log15.NewWriter(log15CtxHandler("pkg", "stdlib", Log15RootHandler)))
	log.SetFlags(0) // don't print a timestamp

	// Set-up the web server routes
	mx := SetupMainMux()
	mx.Handle("/demo/*", demo.NewMux()) // attach demo sub-routes

	// Start the web server
	SetupMainServer("localhost:8080", mx)
}

// log15CtxHandler adds a single context variable to the record being logged
func log15CtxHandler(key string, value interface{}, handler log15.Handler) log15.Handler {
	return log15.FuncHandler(func(r *log15.Record) error {
		r.Ctx = append(r.Ctx, key, value)
		return handler.Log(r)
	})
}

// SetupMainMux returns a goji Mux initialized with the middleware and a health check
// handler. The mux is a net/http.Handler and thus has a ServeHTTP
// method which can be convient to call in tests without the actual network stuff and
// goroutine that SetupMainServer adds
func SetupMainMux() *web.Mux {
	mx := web.New()
	gojiutil.AddCommon15(mx, log15.Root())
	mx.Use(ParamsLogger(log15.Root())) // useful for debugging
	mx.Get("/health-check", healthCheckHandler)
	mx.NotFound(handleNotFound)
	return mx
}

// SetupMainServer allocates a listener socket and starts a web server with graceful restart
// on the specified IP address and port. The ipPort has the format "ip_address:port" or
// ":port" for 0.0.0.0/port.
func SetupMainServer(ipPort string, mux *web.Mux) {
	listener, err := net.Listen("tcp4", ipPort)
	if err != nil {
		FatalError(err.Error())
	}

	// Install our handler at the root of the standard net/http default mux.
	// This allows packages like expvar to continue working as expected.
	mux.Compile()
	http.Handle("/", mux)

	graceful.HandleSignals()
	graceful.PreHook(func() { log15.Warn("Gracefully stopping on signal") })
	graceful.PostHook(func() { log.Printf("Gracefully stopped") })

	err = graceful.Serve(listener, http.DefaultServeMux)
	if err != nil {
		FatalError(err.Error())
	}

	graceful.Wait()
}

// Simple health-check handler that returns the version string
func healthCheckHandler(c web.C, rw http.ResponseWriter, r *http.Request) {
	gojiutil.WriteString(rw, 200, VERSION)
}

// handleNotFound handler, this should be turned into a middleware...
func handleNotFound(c web.C, rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	// probably could print something usefule here...
}

// ParamsLogger logs all query string / form parameters. TODO: move into gojiutils
func ParamsLogger(log15.Logger) web.MiddlewareType {
	return func(c *web.C, h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			params := []interface{}{}
			for k, v := range r.Form {
				params = append(params, k, v[0])
			}
			log15.Debug(r.Method+" "+r.URL.Path, params...)
			//"URLParams", fmt.Sprintf("%+v", c.URLParams))
			//"Env", fmt.Sprintf("%+v", c.Env))
			h.ServeHTTP(rw, r)
		})
	}
}
