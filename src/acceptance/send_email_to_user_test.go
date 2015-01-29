package acceptance

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SendEmailToUser", func() {
	Describe("notification to a user", func() {
		var messageID string

		It("sends a notification to the cf user", func() {
			By("sending the email to the user", func() {
				// SEND A NOTIFICATION TO A USER
				notificationToUserURL := fmt.Sprintf("%s/users/%s", context.NotificationsDomain, context.TestUserGUID)
				output := Run("uaac", "curl", "--insecure", notificationToUserURL, "-X", "POST", "--data", `{"kind_id":"test_notification", "text":"this is a test"}`)

				// VERIFY 200 RESPONSE
				Expect(output).To(ContainSubstring("200 OK"))
				Expect(output).To(ContainSubstring(`"recipient":"` + context.TestUserGUID + `"`))

				// GRAB MESSAGE ID TO CHECK STATUS OF NOTIFICATION
				var notificationResponse NotificationResponse
				jsonBytes := ReturnOnlyBody(output)
				json.Unmarshal(jsonBytes, &notificationResponse)
				Expect(notificationResponse).To(HaveLen(1))
				messageID = notificationResponse[0].ID
			})

			By("verifying the email was sent to the SMTP server", func() {
				notificationToUserStatusURL := fmt.Sprintf("%s/messages/%s", context.NotificationsDomain, messageID)

				Eventually(func() string {
					var statusResponse StatusResponse
					output := Run("uaac", "curl", "--insecure", notificationToUserStatusURL)
					jsonBytes := ReturnOnlyBody(output)
					json.Unmarshal(jsonBytes, &statusResponse)
					return statusResponse.Status
				}, "10s").Should(Equal("delivered"))
			})
		})
	})
})
