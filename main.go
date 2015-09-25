// Copyright (c) 2015 RightScale Inc, All Rights Reserved

package main

import (
	"log"
	"os"

	"github.com/rightscale/uca/log2log15"
	"github.com/rightscale/wstunnel/tunnel"
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
var Log15RootHandler = log15.StreamHandler(os.Stdout, tunnel.SimpleFormat(true))

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

	// Create the internal HTTP servers for the GW and RightLink (really http.Handlers)
	//gwMux := gw.SetupGWMux()

	log15.Info("Entering main loop")
	var c chan struct{}
	<-c // really nothing for us to do, now driven through wstunnels
}

// log15CtxHandler adds a single context variable to the record being logged
func log15CtxHandler(key string, value interface{}, handler log15.Handler) log15.Handler {
	return log15.FuncHandler(func(r *log15.Record) error {
		r.Ctx = append(r.Ctx, key, value)
		return handler.Log(r)
	})
}
