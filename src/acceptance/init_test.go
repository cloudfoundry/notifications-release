package acceptance

import (
	"regexp"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	context      TestSuiteContext
	originalSMTP struct {
		Host string
		Port string
		TLS  string
	}
)

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
		TestClientSenderAuthorities: "notifications.write,emails.write",
		TestClientSenderGrantTypes:  "client_credentials",

		UAACAdminClientID:     LoadOrPanic("UAAC_ADMIN_CLIENT_ID"),
		UAACAdminClientSecret: LoadOrPanic("UAAC_ADMIN_CLIENT_SECRET"),
		CFAdminUsername:       LoadOrPanic("CF_ADMIN_USERNAME"),
		CFAdminPassword:       LoadOrPanic("CF_ADMIN_PASSWORD"),
		NotificationsDomain:   LoadOrPanic("NOTIFICATIONS_DOMAIN"),
		UAADomain:             LoadOrPanic("UAA_DOMAIN"),
		CCDomain:              LoadOrPanic("CC_DOMAIN"),
		NotificationsOrg:      LoadOrPanic("NOTIFICATIONS_ORG"),
		NotificationsSpace:    LoadOrPanic("NOTIFICATIONS_SPACE"),
	}

	// LOGIN AS A CF ADMIN
	Run("cf", "logout")
	Run("cf", "api", context.CCDomain, "--skip-ssl-validation")
	Run("cf", "auth", context.CFAdminUsername, context.CFAdminPassword)

	// PUT NOTIFICATIONS INTO A TESTABLE STATE
	Run("cf", "target", "-o", context.NotificationsOrg, "-s", context.NotificationsSpace)
	Run("cf", "set-env", "notifications", "SMTP_LOGGING_ENABLED", "true")
	Run("cf", "set-env", "notifications", "TRACE", "true")
	Run("cf", "restart", "notifications")
	context.NotificationsAppGUID = strings.TrimSpace(Run("cf", "app", "notifications", "--guid"))

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

	// SET USER AS SPACE DEVELOPER FOR NOTIFICATIONS SPACE
	Run("cf", "set-space-role", context.TestUserName, context.NotificationsOrg, context.NotificationsSpace, "SpaceDeveloper")

	// RETRIEVE THE TEST USER TOKEN
	Run("uaac", "token", "get", context.TestUserName, context.TestUserPassword)
	output = Run("uaac", "context")
	matches := regexp.MustCompile(`access_token: (.*)\n`).FindStringSubmatch(output)
	context.LogToken = matches[1]

	// GET A CLIENT WITH THE RIGHT SCOPES
	Run("uaac", "token", "client", "get", context.UAACAdminClientID, "-s", context.UAACAdminClientSecret)
	Run("uaac", "client", "add", context.TestClientSenderID, "--authorities", context.TestClientSenderAuthorities, "-s", context.TestClientSenderSecret, "--authorized_grant_types", context.TestClientSenderGrantTypes)
})

var _ = AfterSuite(func() {
	// LOGIN AS CF ADMIN
	AlwaysRun("cf", "auth", context.CFAdminUsername, context.CFAdminPassword)
	AlwaysRun("cf", "target", "-o", context.TestOrg, "-s", context.TestSpace)

	// PUT NOTIFICATIONS BACK INTO NORMAL STATE
	Run("cf", "target", "-o", context.NotificationsOrg, "-s", context.NotificationsSpace)
	Run("cf", "unset-env", "notifications", "SMTP_LOGGING_ENABLED")
	Run("cf", "unset-env", "notifications", "TRACE")
	Run("cf", "restart", "notifications")

	// CLEAN UP TEST OBJECTS
	AlwaysRun("cf", "delete-user", context.TestUserName, "-f")
	AlwaysRun("cf", "delete-space", context.TestSpace, "-f")
	AlwaysRun("cf", "delete-org", context.TestOrg, "-f")
	AlwaysRun("uaac", "token", "client", "get", context.UAACAdminClientID, "-s", context.UAACAdminClientSecret)
	AlwaysRun("uaac", "client", "delete", context.TestClientSenderID)
})
