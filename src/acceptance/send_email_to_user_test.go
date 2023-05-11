package acceptance

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"
	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SendEmailToUser", func() {
	var clientToken string

	BeforeEach(func() {
		var err error
		config := warrant.Config{Host: context.UAADomain, SkipVerifySSL: true}
		clientService := warrant.NewClientsService(config)
		clientToken, err = clientService.GetToken(context.TestClientSenderID, context.TestClientSenderSecret)
		Expect(err).NotTo(HaveOccurred())
	})

	It("sends a notification to the cf user", func() {
		var (
			messageID           string
			notificationsClient *support.Client
		)

		By("sending the email to the user", func() {
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
			status, response, err := notificationsClient.Notify.User(clientToken, context.TestUserGUID, notify)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(200))
			Expect(response[0].Recipient).To(Equal(context.TestUserGUID))
			messageID = response[0].NotificationID
		})

		By("verifying the email was sent to the SMTP server", func() {
			Eventually(func() string {
				_, message, err := notificationsClient.Messages.Get(clientToken, messageID)
				Expect(err).NotTo(HaveOccurred())
				return message.Status
			}, 10*time.Second).Should(Equal("delivered"))
		})
	})
})
