package acceptance

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"
	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SendEmailToSpace", func() {
	var clientToken string

	BeforeEach(func() {
		var err error

		config := warrant.Config{Host: context.UAADomain, SkipVerifySSL: true}
		clientService := warrant.NewClientsService(config)
		clientToken, err = clientService.GetToken(context.TestClientSenderID, context.TestClientSenderSecret)
		Expect(err).NotTo(HaveOccurred())
	})

	It("sends a notification a user within a particular space", func() {
		var (
			messageID           string
			notificationsClient *support.Client
		)

		By("sending the notification to the space", func() {
			// PUT A USER IN A SPACE
			Run("cf", "auth", context.CFAdminUsername, context.CFAdminPassword)
			Run("cf", "set-space-role", context.TestUserName, context.TestOrg, context.TestSpace, "SpaceDeveloper")
			Run("cf", "auth", context.TestUserName, context.TestUserPassword)

			// GRAB IDs FOR CURL REQUESTS
			output := Run("cf", "curl", fmt.Sprintf("/v2/spaces?q=name:%s", context.TestSpace))

			var spaceGUIDResponse GUIDResponse
			err := json.Unmarshal([]byte(output), &spaceGUIDResponse)
			if err != nil {
				panic(err)
			}
			spaceGUID := spaceGUIDResponse.Resources[0].MetaData.GUID

			notificationsClient = support.NewClient(context.NotificationsDomain)
			transport := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			notificationsClient.HTTPClient = &http.Client{Transport: transport}
			notify := support.Notify{
				KindID: "test_notification",
				Text:   "this is a test",
				HTML:   "something",
			}
			status, responses, err := notificationsClient.Notify.Space(clientToken, spaceGUID, notify)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(200))
			Expect(responses).To(HaveLen(2))

			for _, response := range responses {
				if response.Recipient == context.TestUserGUID {
					messageID = response.NotificationID
				}
			}

			Expect(messageID).NotTo(BeEmpty())
		})

		By("verifying the notification was sent to the space user", func() {
			Eventually(func() string {
				_, message, err := notificationsClient.Messages.Get(clientToken, messageID)
				Expect(err).NotTo(HaveOccurred())
				return message.Status
			}, 10*time.Second).Should(Equal("delivered"))
		})
	})
})
