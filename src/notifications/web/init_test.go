package web_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWebSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "web")
}
