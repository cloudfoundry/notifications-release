package warrant_test

import (
	"io"
	"os"
	"testing"

	"github.com/pivotal-cf-experimental/warrant/internal/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	fakeUAAServer *fakes.UAAServer
	TraceWriter   io.Writer
)

func TestWarrantSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Warrant Suite")
}

var _ = BeforeSuite(func() {
	if os.Getenv("TRACE") == "true" {
		TraceWriter = os.Stdout
	}

	fakeUAAServer = fakes.NewUAAServer()
	fakeUAAServer.Start()
})

var _ = AfterSuite(func() {
	fakeUAAServer.Close()
})

var _ = BeforeEach(func() {
	fakeUAAServer.Reset()
})
