// Copyright (c) 2015 RightScale, Inc. - see LICENSE

package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/rightscale/uca/abc"
	"gopkg.in/inconshreveable/log15.v2"
)

func init() {
	abc.InTest = true
}

func TestUCA(t *testing.T) {
	//log.SetOutput(ginkgo.GinkgoWriter)
	Log15RootHandler = log15.StreamHandler(GinkgoWriter, log15.TerminalFormat())
	log15.Root().SetHandler(Log15RootHandler)

	format.UseStringerRepresentation = true
	RegisterFailHandler(Fail)

	RunSpecs(t, "UCA")
}
