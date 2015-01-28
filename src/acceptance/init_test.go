package acceptance

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var context TestSuiteContext

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var _ = BeforeSuite(func() {
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
		LoggregatorDomain:     LoadOrPanic("LOGGREGATOR_DOMAIN"),
	}

	// LOGIN AS A CF USER
	Run("cf", "api", context.CCDomain, "--skip-ssl-validation")
	Run("cf", "login", "-u", context.CFAdminUsername, "-p", context.CFAdminPassword)

	// PUT NOTIFICATIONS INTO A TESTABLE STATE
	Run("cf", "target", "-o", context.NotificationsOrg, "-s", context.NotificationsSpace)
	Run("cf", "set-env", "notifications", "SMTP_LOGGING_ENABLED", "true")
	Run("cf", "restart", "notifications")
	context.NotificationsAppGUID = Run("cf", "app", "notifications", "--guid")

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

var _ = AfterSuite(func() {
	AlwaysRun("cf", "login", "-u", context.CFAdminUsername, "-p", context.CFAdminPassword)
	AlwaysRun("cf", "target", "-o", context.TestOrg, "-s", context.TestSpace)
	AlwaysRun("cf", "delete-user", context.TestUserName, "-f")
	AlwaysRun("cf", "delete-space", context.TestSpace, "-f")
	AlwaysRun("cf", "delete-org", context.TestOrg, "-f")
	AlwaysRun("uaac", "token", "client", "get", context.UAACAdminClientID, "-s", context.UAACAdminClientSecret)
	AlwaysRun("uaac", "client", "delete", context.TestClientSenderID)
})
