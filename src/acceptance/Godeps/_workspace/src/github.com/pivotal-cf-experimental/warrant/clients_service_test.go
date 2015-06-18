package warrant_test

import (
	"time"

	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientsService", func() {
	var (
		service warrant.ClientsService
		token   string
		config  warrant.Config
	)

	BeforeEach(func() {
		config = warrant.Config{
			Host:          fakeUAAServer.URL(),
			SkipVerifySSL: true,
			TraceWriter:   TraceWriter,
		}
		service = warrant.NewClientsService(config)
		token = fakeUAAServer.ClientTokenFor("admin", []string{"clients.write", "clients.read"}, []string{"clients"})
	})

	Describe("Create/Get", func() {
		var client warrant.Client

		BeforeEach(func() {
			client = warrant.Client{
				ID:                   "client-id",
				Scope:                []string{"openid"},
				ResourceIDs:          []string{"none"},
				Authorities:          []string{"scim.read", "scim.write"},
				AuthorizedGrantTypes: []string{"client_credentials"},
				AccessTokenValidity:  5000 * time.Second,
			}
		})

		It("an error does not occur and the new client can be fetched", func() {
			err := service.Create(client, "client-secret", token)
			Expect(err).NotTo(HaveOccurred())

			foundClient, err := service.Get(client.ID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(foundClient).To(Equal(client))
		})

		It("responds with an error when the client cannot be created", func() {
			client.AuthorizedGrantTypes = []string{"invalid-grant-type"}
			err := service.Create(client, "client-secret", token)
			Expect(err).To(BeAssignableToTypeOf(warrant.UnexpectedStatusError{}))
		})

		It("responds with an error when the client cannot be found", func() {
			_, err := service.Get("unknown-client", token)
			Expect(err).To(BeAssignableToTypeOf(warrant.NotFoundError{}))
		})
	})

	Describe("GetToken", func() {
		var (
			client       warrant.Client
			clientSecret string
		)

		BeforeEach(func() {
			client = warrant.Client{
				ID:                   "client-id",
				Scope:                []string{"openid"},
				ResourceIDs:          []string{"none"},
				Authorities:          []string{"scim.read", "scim.write"},
				AuthorizedGrantTypes: []string{"client_credentials"},
				AccessTokenValidity:  5000 * time.Second,
			}
			clientSecret = "client-secret"

			err := service.Create(client, clientSecret, token)
			Expect(err).NotTo(HaveOccurred())
		})

		It("retrieves a token for the client given a valid secret", func() {
			clientToken, err := service.GetToken(client.ID, clientSecret)
			Expect(err).NotTo(HaveOccurred())

			tokensService := warrant.NewTokensService(config)
			decodedToken, err := tokensService.Decode(clientToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(decodedToken.ClientID).To(Equal(client.ID))
		})
	})

	Describe("Delete", func() {
		var client warrant.Client

		BeforeEach(func() {
			client = warrant.Client{
				ID:                   "client-id",
				Scope:                []string{"openid"},
				ResourceIDs:          []string{"none"},
				Authorities:          []string{"scim.read", "scim.write"},
				AuthorizedGrantTypes: []string{"client_credentials"},
				AccessTokenValidity:  5000 * time.Second,
			}

			err := service.Create(client, "secret", token)
			Expect(err).NotTo(HaveOccurred())
		})

		It("deletes the client", func() {
			err := service.Delete(client.ID, token)
			Expect(err).NotTo(HaveOccurred())

			_, err = service.Get(client.ID, token)
			Expect(err).To(BeAssignableToTypeOf(warrant.NotFoundError{}))
		})

		It("errors when the token is unauthorized", func() {
			token = fakeUAAServer.ClientTokenFor("admin", []string{"clients.foo", "clients.boo"}, []string{"clients"})
			err := service.Delete(client.ID, token)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(warrant.UnauthorizedError{}))
		})
	})
})
