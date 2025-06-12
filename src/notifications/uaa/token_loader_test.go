package uaa_test

import (
	"github.com/cloudfoundry/notifications-release/src/notifications/v81/testing/mocks"
	"github.com/cloudfoundry/notifications-release/src/notifications/v81/uaa"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TokenLoader", func() {
	Describe("#Load", func() {
		It("Gets a zoned client token based on hostname", func() {
			uaaClient := mocks.NewZonedUAAClient()
			uaaClient.GetClientTokenCall.Returns.Token = "my-fake-token"

			tokenLoader := uaa.NewTokenLoader(uaaClient)

			token, err := tokenLoader.Load("my-uaa-zone")
			Expect(token).To(Equal("my-fake-token"))
			Expect(err).To(BeNil())

			Expect(uaaClient.GetClientTokenCall.Receives.Host).To(Equal("my-uaa-zone"))
		})
	})
})
