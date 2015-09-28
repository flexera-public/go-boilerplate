// Copyright (c) 2015 RightScale, Inc. - see LICENSE

package demo

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"gopkg.in/inconshreveable/log15.v2"
)

func TestDemo(t *testing.T) {
	// send the logs through the GinkgoWriter, which buffers up the output for each test
	// and only prints it if the test fails. Use ginkgo -v to always see the output.
	log15RootHandler := log15.StreamHandler(GinkgoWriter, log15.TerminalFormat())
	log15.Root().SetHandler(log15RootHandler)

	format.UseStringerRepresentation = true
	RegisterFailHandler(Fail)

	RunSpecs(t, "Demo")
}
