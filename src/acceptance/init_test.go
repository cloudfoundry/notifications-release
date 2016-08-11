package acceptance

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	context              TestSuiteContext
	environmentVariables map[string]string
	originalSMTP         struct {
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
	context = TestSuiteContext{
		TestUserName:                Randomized("user"),
		TestUserPassword:            Randomized("password"),
		TestOrg:                     Randomized("org"),
		TestSpace:                   Randomized("space"),
		TestClientSenderID:          Randomized("client"),
		TestClientSenderSecret:      Randomized("secret"),
		TestClientSenderAuthorities: []string{"notifications.write", "emails.write"},
		TestClientSenderGrantTypes:  []string{"client_credentials"},

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
	Run("cf", "target", "-o", context.NotificationsOrg, "-s", context.NotificationsSpace)

	saveNotificationsEnvironmentVariables()

	// PUT NOTIFICATIONS INTO A TESTABLE STATE
	Run("cf", "set-env", "notifications", "SMTP_LOGGING_ENABLED", "true")
	Run("cf", "set-env", "notifications", "TRACE", "true")
	Run("cf", "restart", "notifications")
	context.NotificationsAppGUID = strings.TrimSpace(Run("cf", "app", "notifications", "--guid"))

	// CREATE A USER AND GRAB ITS TOKEN
	Run("cf", "create-user", context.TestUserName, context.TestUserPassword)
	Run("cf", "create-org", context.TestOrg)
	Run("cf", "create-space", context.TestSpace, "-o", context.TestOrg)
	Run("cf", "target", "-o", context.TestOrg, "-s", context.TestSpace)

	setupTestUser()

	// SET USER AS SPACE DEVELOPER FOR NOTIFICATIONS SPACE
	Run("cf", "set-space-role", context.TestUserName, context.NotificationsOrg, context.NotificationsSpace, "SpaceDeveloper")
})

var _ = AfterSuite(func() {
	// LOGIN AS CF ADMIN
	AlwaysRun("cf", "auth", context.CFAdminUsername, context.CFAdminPassword)
	AlwaysRun("cf", "target", "-o", context.TestOrg, "-s", context.TestSpace)

	// PUT NOTIFICATIONS BACK INTO NORMAL STATE
	Run("cf", "target", "-o", context.NotificationsOrg, "-s", context.NotificationsSpace)
	Run("cf", "unset-env", "notifications", "SMTP_LOGGING_ENABLED")
	Run("cf", "unset-env", "notifications", "TRACE")

	restoreNotificationsEnvironmentVariables()
	Run("cf", "restart", "notifications")

	// CLEAN UP TEST OBJECTS
	AlwaysRun("cf", "delete-user", context.TestUserName, "-f")
	AlwaysRun("cf", "delete-space", context.TestSpace, "-f")
	AlwaysRun("cf", "delete-org", context.TestOrg, "-f")

	teardownTestUser()
})

func saveNotificationsEnvironmentVariables() {
	guid := Run("cf", "app", "notifications", "--guid")
	environmentJSON := Run("cf", "curl", fmt.Sprintf("/v2/apps/%s/env", guid))
	var env struct {
		EnvironmentJSON map[string]string `json:"environment_json"`
	}
	err := json.Unmarshal([]byte(environmentJSON), &env)
	Expect(err).NotTo(HaveOccurred())
	environmentVariables = env.EnvironmentJSON
}

func restoreNotificationsEnvironmentVariables() {
	for name, value := range environmentVariables {
		Run("cf", "set-env", "notifications", name, value)
	}
}

func setupTestUser() {
	config := warrant.Config{Host: context.UAADomain, SkipVerifySSL: true}
	adminToken := fetchAdminToken(config)

	userService := warrant.NewUsersService(config)
	users, err := userService.List(warrant.Query{Filter: fmt.Sprintf("username eq '%s'", context.TestUserName)}, adminToken)
	Expect(err).NotTo(HaveOccurred())
	testUser := users[0]

	testUser.Emails = []string{"this-is-an-example@example.com"}

	userService.Update(testUser, adminToken)
	context.TestUserGUID = testUser.ID

	context.LogToken, err = userService.GetToken(context.TestUserName, context.TestUserPassword, warrant.Client{
		ID: "cf",
		Scope: []string{
			"cloud_controller.read",
			"cloud_controller.write",
			"openid",
			"password.write",
			"cloud_controller.admin",
			"cloud_controller.admin_read_only",
			"scim.read",
			"scim.write",
			"doppler.firehose",
			"uaa.user",
			"routing.router_groups.read",
			"routing.router_groups.write",
		},
	})
	Expect(err).NotTo(HaveOccurred())

	clientService := warrant.NewClientsService(config)
	client := warrant.Client{
		ID:                   context.TestClientSenderID,
		Authorities:          context.TestClientSenderAuthorities,
		AuthorizedGrantTypes: context.TestClientSenderGrantTypes,
		AccessTokenValidity:  time.Duration(6000 * time.Second),
	}
	Expect(clientService.Create(client, context.TestClientSenderSecret, adminToken)).To(Succeed())
}

func teardownTestUser() {
	config := warrant.Config{Host: context.UAADomain, SkipVerifySSL: true}
	adminToken := fetchAdminToken(config)

	clientService := warrant.NewClientsService(config)
	Expect(clientService.Delete(context.TestClientSenderID, adminToken)).To(Succeed())
}

func fetchAdminToken(config warrant.Config) string {
	clientService := warrant.NewClientsService(config)

	adminToken, err := clientService.GetToken(context.UAACAdminClientID, context.UAACAdminClientSecret)
	Expect(err).NotTo(HaveOccurred())

	return adminToken
}
