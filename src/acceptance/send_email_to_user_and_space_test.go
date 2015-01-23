package acceptance

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SendEmailToUserAndSpace", func() {
	var context TestSuiteContext

	BeforeSuite(func() {
		// LOAD ENVIRONMENT AND CONTEXT
		context = TestSuiteContext{
			TestUserName:                Randomized("user"),
			TestUserPassword:            Randomized("password"),
			TestOrg:                     Randomized("org"),
			TestSpace:                   Randomized("space"),
			TestClientSenderID:          Randomized("client"),
			TestClientSenderSecret:      Randomized("secret"),
			TestClientSenderAuthorities: "notifications.write",
			TestClientSenderGrantTypes:  "client_credentials",

			UAACAdminClientID:     LoadOrPanic("UAAC_ADMIN_CLIENT_ID"),
			UAACAdminClientSecret: LoadOrPanic("UAAC_ADMIN_CLIENT_SECRET"),
			CFAdminUsername:       LoadOrPanic("CF_ADMIN_USERNAME"),
			CFAdminPassword:       LoadOrPanic("CF_ADMIN_PASSWORD"),
			NotificationsDomain:   LoadOrPanic("NOTIFICATIONS_DOMAIN"),
			UAADomain:             LoadOrPanic("UAA_DOMAIN"),
			CCDomain:              LoadOrPanic("CC_DOMAIN"),
		}

		// LOGIN AS A CF USER
		Run("cf", "api", context.CCDomain, "--skip-ssl-validation")
		Run("cf", "login", "-u", context.CFAdminUsername, "-p", context.CFAdminPassword)

		// CREATE A USER AND GRAB ITS TOKEN
		Run("cf", "create-user", context.TestUserName, context.TestUserPassword)
		Run("cf", "create-org", context.TestOrg)
		Run("cf", "create-space", context.TestSpace, "-o", context.TestOrg)
		Run("cf", "target", "-o", context.TestOrg, "-s", context.TestSpace)
		Run("uaac", "target", context.UAADomain, "--skip-ssl-validation")
		Run("uaac", "token", "client", "get", context.UAACAdminClientID, "-s", context.UAACAdminClientSecret)
		Run("uaac", "user", "update", context.TestUserName, "--emails", "this-is-an-example@example.com")

		output := Run("uaac", "user", "get", context.TestUserName, "-a", "id")
		context.TestUserGUID = strings.TrimSpace(strings.Split(output, ":")[1])

		// GET A CLIENT WITH THE RIGHT SCOPES
		Run("uaac", "client", "add", context.TestClientSenderID, "--authorities", context.TestClientSenderAuthorities, "-s", context.TestClientSenderSecret, "--authorized_grant_types", context.TestClientSenderGrantTypes)
		Run("uaac", "token", "client", "get", context.TestClientSenderID, "-s", context.TestClientSenderSecret)
	})

	AfterSuite(func() {
		AlwaysRun("cf", "login", "-u", context.CFAdminUsername, "-p", context.CFAdminPassword)
		AlwaysRun("cf", "target", "-o", context.TestOrg, "-s", context.TestSpace)
		AlwaysRun("cf", "delete-user", context.TestUserName, "-f")
		AlwaysRun("cf", "delete-space", context.TestSpace, "-f")
		AlwaysRun("cf", "delete-org", context.TestOrg, "-f")
		AlwaysRun("uaac", "token", "client", "get", context.UAACAdminClientID, "-s", context.UAACAdminClientSecret)
		AlwaysRun("uaac", "client", "delete", context.TestClientSenderID)
	})

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
				}, 10*time.Second).Should(Equal("delivered"))
			})
		})
	})

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
