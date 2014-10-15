package acceptance

import (
    "encoding/json"
    "fmt"
    "strings"

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

            UAACAdminClientID:       LoadOrPanic("UAAC_ADMIN_CLIENT_USERNAME"),
            UAACAdminClientPassword: LoadOrPanic("UAAC_ADMIN_CLIENT_PASSWORD"),
            CFAdminUsername:         LoadOrPanic("CF_ADMIN_USERNAME"),
            CFAdminPassword:         LoadOrPanic("CF_ADMIN_PASSWORD"),
            NotificationsDomain:     LoadOrPanic("NOTIFICATIONS_DOMAIN"),
            UAADomain:               LoadOrPanic("UAA_DOMAIN"),
            CCDomain:                LoadOrPanic("CC_DOMAIN"),
        }

        // LOGIN AS A CF USER
        Run("cf", "api", context.CCDomain, "--skip-ssl-validation")
        Run("cf", "login", "-u", context.CFAdminUsername, "-p", context.CFAdminPassword)

        // CREATE A USER AND GRAB ITS TOKEN
        Run("cf", "create-user", context.TestUserName, context.TestUserPassword)
        Run("cf", "create-org", context.TestOrg)
        Run("cf", "create-space", context.TestSpace, "-o", context.TestOrg)
        Run("cf", "target", "-o", context.TestOrg, "-s", context.TestSpace)
        Run("uaac", "target", context.UAADomain)
        Run("uaac", "token", "client", "get", context.UAACAdminClientID, "-s", context.UAACAdminClientPassword)

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
        AlwaysRun("uaac", "token", "client", "get", context.UAACAdminClientID, "-s", context.UAACAdminClientPassword)
        AlwaysRun("uaac", "client", "delete", context.TestClientSenderID)
    })

    Describe("sending an email to a user", func() {
        It("returns a 200", func() {
            // SEND A NOTIFICATION TO A USER
            notificationToUserURL := fmt.Sprintf("%s/users/%s", context.NotificationsDomain, context.TestUserGUID)
            output := Run("uaac", "curl", notificationToUserURL, "-X", "POST", "--data", `{"kind_id":"test_notification", "text":"this is a test"}`)

            // VERIFY 200 RESPONSE
            Expect(output).To(ContainSubstring("200 OK"))
            Expect(output).To(ContainSubstring(`"recipient":"` + context.TestUserGUID + `"`))
        })
    })

    Describe("sending an email to a space", func() {
        It("returns a 200", func() {
            // PUT A USER IN A SPACE
            Run("cf", "set-space-role", context.TestUserName, context.TestOrg, context.TestSpace, "SpaceDeveloper")
            Run("cf", "login", "-u", context.TestUserName, "-p", context.TestUserPassword)

            // GRAB IDs FOR CURL REQUESTS
            output := Run("cf", "curl", "/v2/organizations")

            var orgGUIDResponse GUIDResponse
            err := json.Unmarshal([]byte(output), &orgGUIDResponse)
            if err != nil {
                panic(err)
            }
            orgGUID := orgGUIDResponse.Resources[0].MetaData.GUID

            output = Run("cf", "curl", fmt.Sprintf("/v2/organizations/%s/spaces", orgGUID))

            var spaceGUIDResponse GUIDResponse
            err = json.Unmarshal([]byte(output), &spaceGUIDResponse)
            if err != nil {
                panic(err)
            }
            spaceGUID := spaceGUIDResponse.Resources[0].MetaData.GUID

            // SEND A NOTIFICATION TO A SPACE
            notificationToSpaceURL := fmt.Sprintf("%s/spaces/%s", context.NotificationsDomain, spaceGUID)
            output = Run("uaac", "curl", notificationToSpaceURL, "-X", "POST", "--data", `{"kind_id":"test_notification", "text":"this is a test"}`)

            // VERIFY 200 RESPONSE
            Expect(output).To(ContainSubstring("200 OK"))
            Expect(output).To(ContainSubstring(`"recipient":"` + context.TestUserGUID + `"`))
        })
    })
})
