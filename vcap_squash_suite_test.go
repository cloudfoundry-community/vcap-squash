package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestVcapSquash(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VcapSquash Suite")
}
