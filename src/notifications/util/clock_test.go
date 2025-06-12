package util_test

import (
	"time"

	"github.com/cloudfoundry/notifications-release/src/notifications/v81/util"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clock", func() {
	Describe("Now", func() {
		It("should return the current time", func() {
			clock := util.NewClock()

			currentTime := clock.Now()
			Expect(currentTime).To(BeTemporally("~", time.Now()))
		})
	})
})
