// +build performance

package acceptance

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/chrj/smtpd"
	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	userCount   = 400
	workerCount = 10

	// This is a huge hack and should be fixed at some point. The admin email
	// account does not have a valid email address and so will not produce a
	// deliverable email.
	undeliverableEmailCount = 2
)

type User struct {
	Name    string
	OrgRole string
	GUID    string
}

type UserCreator struct {
	In      <-chan User
	Out     chan<- User
	OrgName string
}

func (uc UserCreator) Run() {
	for {
		user := <-uc.In
		Run("cf", "create-user", user.Name, "password")
		Run("cf", "set-space-role", user.Name, uc.OrgName, "benchmark", "SpaceDeveloper")
		Run("cf", "set-org-role", user.Name, uc.OrgName, user.OrgRole)

		config := warrant.Config{Host: context.UAADomain, SkipVerifySSL: true}
		adminToken := fetchAdminToken(config)
		userService := warrant.NewUsersService(config)
		users, err := userService.Find(warrant.UsersQuery{Filter: fmt.Sprintf("username eq '%s'", user.Name)}, adminToken)
		Expect(err).NotTo(HaveOccurred())
		userGUID := users[0].ID

		user.GUID = userGUID
		uc.Out <- user
	}
}

type Message struct {
	URL      string
	Payload  map[string]interface{}
	Response *http.Response
}

type MessageSender struct {
	Client *http.Client
	Token  string
	In     <-chan Message
	Out    chan<- Message
}

func (ms MessageSender) Run() {
	for {
		message := <-ms.In
		content, err := json.Marshal(message.Payload)
		if err != nil {
			panic(err)
		}

		request, err := http.NewRequest("POST", message.URL, bytes.NewBuffer(content))
		if err != nil {
			panic(err)
		}

		request.Header.Set("Authorization", "Bearer "+ms.Token)

		message.Response, err = ms.Client.Do(request)
		if err != nil {
			panic(err)
		}

		ms.Out <- message
	}
}

var _ = Describe("Performance", func() {
	var (
		orgName    string
		orgGUID    string
		spaceGUID  string
		users      []User
		token      string
		messageIDs []string
	)

	BeforeEach(func() {
		smtpPort := freePort()
		smtpHost := LoadOrPanic("SMTP_HOST")
		smtpServer := &smtpd.Server{
			Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {
				context.Deliveries = append(context.Deliveries, env)
				return nil
			},
		}
		go smtpServer.ListenAndServe(fmt.Sprintf("%s:%s", smtpHost, smtpPort))

		context.Deliveries = []smtpd.Envelope{}
		Run("cf", "auth", context.CFAdminUsername, context.CFAdminPassword)
		Run("cf", "target", "-o", context.NotificationsOrg, "-s", context.NotificationsSpace)

		originalEnv := Run("cf", "env", "notifications")
		originalSMTP.Host = regexp.MustCompile(`SMTP_HOST:\s(.+)\n`).FindStringSubmatch(originalEnv)[1]
		originalSMTP.Port = regexp.MustCompile(`SMTP_PORT:\s(.+)\n`).FindStringSubmatch(originalEnv)[1]
		originalSMTP.TLS = regexp.MustCompile(`SMTP_TLS:\s(.+)\n`).FindStringSubmatch(originalEnv)[1]

		Run("cf", "set-env", "notifications", "SMTP_HOST", smtpHost)
		Run("cf", "set-env", "notifications", "SMTP_PORT", smtpPort)
		Run("cf", "set-env", "notifications", "SMTP_TLS", "false")
		Run("cf", "restart", "notifications")

		orgName = Randomized("org")
		Run("cf", "create-org", orgName)
		Run("cf", "target", "-o", orgName)
		Run("cf", "create-space", "benchmark")
		Run("cf", "target", "-s", "benchmark")

		orgGUID = strings.TrimSpace(Run("cf", "org", orgName, "--guid"))
		spaceGUID = strings.TrimSpace(Run("cf", "space", "benchmark", "--guid"))
		users = []User{}

		userInChan := make(chan User, userCount)
		userOutChan := make(chan User)

		for i := 0; i < userCount; i++ {
			username := fmt.Sprintf("user-%d@%s.example.com", i, orgName)
			var orgRole string
			switch i % 3 {
			case 0:
				orgRole = "OrgManager"
			case 1:
				orgRole = "BillingManager"
			case 2:
				orgRole = "OrgAuditor"
			}
			userInChan <- User{Name: username, OrgRole: orgRole}
		}

		for i := 0; i < workerCount; i++ {
			uc := UserCreator{
				In:      userInChan,
				Out:     userOutChan,
				OrgName: orgName,
			}
			go uc.Run()
		}

		for i := 0; i < userCount; i++ {
			var user User
			Eventually(userOutChan, 60*time.Second).Should(Receive(&user))

			users = append(users, user)
		}

		clientService := warrant.NewClientsService(warrant.Config{Host: context.UAADomain, SkipVerifySSL: true})
		var err error
		token, err = clientService.GetToken(context.TestClientSenderID, context.TestClientSenderSecret)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Run("cf", "target", "-o", context.NotificationsOrg, "-s", context.NotificationsSpace)

		if originalSMTP.Host != "" {
			Run("cf", "set-env", "notifications", "SMTP_HOST", originalSMTP.Host)
			Run("cf", "set-env", "notifications", "SMTP_PORT", originalSMTP.Port)
			Run("cf", "set-env", "notifications", "SMTP_TLS", originalSMTP.TLS)
			Run("cf", "restart", "notifications")
		}

		Run("cf", "delete-org", orgName, "-f")
		Run("cf", "logout")
	})

	It("sending email to 500 users", func() {
		var messages []Message

		By("creating a notification for the space", func() {
			messages = append(messages, Message{
				URL: fmt.Sprintf("%s/spaces/%s", context.NotificationsDomain, spaceGUID),
				Payload: map[string]interface{}{
					"kind_id": "test_notification",
					"text":    "this is a test",
					"html":    "something",
				},
			})
		})

		By("creating a notification for the organization", func() {
			messages = append(messages, Message{
				URL: fmt.Sprintf("%s/organizations/%s", context.NotificationsDomain, orgGUID),
				Payload: map[string]interface{}{
					"kind_id": "test_notification",
					"text":    "this is a test",
					"html":    "something",
				},
			})
		})

		By("creating a notification for each role in the organization", func() {
			for _, role := range []string{"OrgManager", "BillingManager", "OrgAuditor"} {
				messages = append(messages, Message{
					URL: fmt.Sprintf("%s/organizations/%s", context.NotificationsDomain, orgGUID),
					Payload: map[string]interface{}{
						"role":    role,
						"kind_id": "test_notification",
						"text":    "this is a test",
						"html":    "something",
					},
				})
			}
		})

		By("creating a notification for each user", func() {
			for _, u := range users {
				messages = append(messages, Message{
					URL: fmt.Sprintf("%s/users/%s", context.NotificationsDomain, u.GUID),
					Payload: map[string]interface{}{
						"kind_id": "test_notification",
						"text":    "this is a test",
						"html":    "something",
					},
				})
			}
		})

		By("creating a notification for each email address", func() {
			for _, u := range users {
				messages = append(messages, Message{
					URL: fmt.Sprintf("%s/emails", context.NotificationsDomain),
					Payload: map[string]interface{}{
						"to":      u.Name,
						"kind_id": "test_notification",
						"text":    "this is a test",
						"html":    "something",
					},
				})
			}
		})

		messageInChan := make(chan Message, len(messages))
		messageOutChan := make(chan Message)

		By("setting up some workers to send notifications", func() {
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}

			for i := 0; i < workerCount; i++ {
				ms := MessageSender{
					Client: client,
					Token:  token,
					In:     messageInChan,
					Out:    messageOutChan,
				}
				go ms.Run()
			}

			for _, message := range messages {
				messageInChan <- message
			}
		})

		By("waiting for all notifications to be sent", func() {
			for i := 0; i < len(messages); i++ {
				var message Message
				Eventually(messageOutChan, 5*time.Minute).Should(Receive(&message))

				if message.Response.StatusCode == http.StatusOK {
					var responses []struct {
						MessageID string `json:"message_id"`
						Recipient string `json:"recipient"`
						Status    string `json:"status"`
					}

					err := json.NewDecoder(message.Response.Body).Decode(&responses)
					Expect(err).NotTo(HaveOccurred())

					for _, response := range responses {
						messageIDs = append(messageIDs, response.MessageID)
					}
				} else {
					body, _ := ioutil.ReadAll(message.Response.Body)
					Fail(fmt.Sprintf("Received bad response (%d) %s", message.Response.StatusCode, body))
				}
			}
		})

		By("waiting to receive emails for each recipient", func() {
			expectedDeliveryCount := len(messageIDs) - undeliverableEmailCount
			Eventually(func() int { return len(context.Deliveries) }, 5*time.Minute).Should(Equal(expectedDeliveryCount))
		})
	})
})
