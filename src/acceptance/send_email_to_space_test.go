package acceptance

import (
	"encoding/json"
	"fmt"
	"time"
)

var _ = Describe("SendEmailToSpace", func() {
	Describe("notification to a space", func() {
		var messageID string

		It("sends a notification a user within a particular space", func() {
			By("sending the notification to the space", func() {
				// PUT A USER IN A SPACE
				Run("cf", "set-space-role", context.TestUserName, context.TestOrg, context.TestSpace, "SpaceDeveloper")
				Run("cf", "login", "-u", context.TestUserName, "-p", context.TestUserPassword)

				// GRAB IDs FOR CURL REQUESTS
				output := Run("cf", "curl", fmt.Sprintf("/v2/spaces?q=name:%s", context.TestSpace))

				var spaceGUIDResponse GUIDResponse
				err := json.Unmarshal([]byte(output), &spaceGUIDResponse)
				if err != nil {
					panic(err)
				}
				spaceGUID := spaceGUIDResponse.Resources[0].MetaData.GUID

				// SEND A NOTIFICATION TO A SPACE
				notificationToSpaceURL := fmt.Sprintf("%s/spaces/%s", context.NotificationsDomain, spaceGUID)
				output = Run("uaac", "curl", "--insecure", notificationToSpaceURL, "-X", "POST", "--data", `{"kind_id":"test_notification", "text":"this is a test"}`)

				// VERIFY 200 RESPONSE
				Expect(output).To(ContainSubstring("200 OK"))
				Expect(output).To(ContainSubstring(`"recipient":"` + context.TestUserGUID + `"`))

				// GRAB MESSAGE ID TO CHECK STATUS OF NOTIFICATION
				var notificationResponses NotificationResponse
				jsonBytes := ReturnOnlyBody(output)
				json.Unmarshal(jsonBytes, &notificationResponses)
				Expect(notificationResponses).To(HaveLen(2))

				for _, response := range notificationResponses {
					if response.Recipient == context.TestUserGUID {
						messageID = response.ID
					}
				}
			})

			By("verifying the notification was sent to the space user", func() {
				notificationToUserStatusURL := fmt.Sprintf("%s/messages/%s", context.NotificationsDomain, messageID)

				Eventually(func() string {
					var statusResponse StatusResponse
					output := Run("uaac", "curl", "--insecure", notificationToUserStatusURL)
					jsonBytes := ReturnOnlyBody(output)
					json.Unmarshal(jsonBytes, &statusResponse)
					return statusResponse.Status
				}, 10*time.Second).Should(Equal("delivered"))
			})
		})
	})
})
