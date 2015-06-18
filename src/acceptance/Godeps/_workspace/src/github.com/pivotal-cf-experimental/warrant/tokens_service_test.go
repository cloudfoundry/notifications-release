package warrant_test

import (
	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TokensService", func() {
	var service warrant.TokensService

	BeforeEach(func() {
		service = warrant.NewTokensService(warrant.Config{
			Host:          fakeUAAServer.URL(),
			SkipVerifySSL: true,
			TraceWriter:   TraceWriter,
		})
	})

	Describe("Decode", func() {
		var encodedToken string

		BeforeEach(func() {
			encodedToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci1pZCIsInNjb3BlIjpbInNjaW0ucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIuYWRtaW4iLCJwYXNzd29yZC53cml0ZSJdfQ.QWNTRFahfCn7qSWxEHTCn6QeZMJxNMq9a_TP8aANc4k"
		})

		It("returns a decoded token given an encoded token string", func() {
			token, err := service.Decode(encodedToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal(warrant.Token{
				UserID: "user-id",
				Scopes: []string{
					"scim.read",
					"cloud_controller.admin",
					"password.write",
				},
			}))
		})
	})
})
