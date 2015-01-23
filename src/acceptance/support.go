package acceptance

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nu7hatch/gouuid"
)

type TestSuiteContext struct {
	UAACAdminClientID           string
	UAACAdminClientSecret       string
	CFAdminUsername             string
	CFAdminPassword             string
	TestUserGUID                string
	TestUserName                string
	TestUserPassword            string
	TestOrg                     string
	TestSpace                   string
	TestClientSenderID          string
	TestClientSenderSecret      string
	TestClientSenderAuthorities string
	TestClientSenderGrantTypes  string
	NotificationsDomain         string
	UAADomain                   string
	CCDomain                    string
}

type GUIDResponse struct {
	Resources []struct {
		MetaData struct {
			GUID string `json:"guid"`
		} `json:"metadata"`
	}
}

type NotificationResponse []struct {
	ID        string `json:"notification_id"`
	Recipient string `json:"recipient"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func AlwaysRun(command string, arguments ...string) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}()

	Run(command, arguments...)
}

func Run(command string, arguments ...string) string {
	parts := []string{"$", command}
	parts = append(parts, arguments...)
	fmt.Println(strings.Join(parts, " "))

	cmd := exec.Command(command, arguments...)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		panic(err)
	}

	return string(output)
}

func LoadOrPanic(variable string) string {
	value := os.Getenv(variable)
	if value == "" {
		panic(variable + " is a required environment variable")
	}
	return value
}

func Randomized(prefix string) string {
	entropy, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return prefix + "-" + entropy.String()
}

func ReturnOnlyBody(body string) []byte {
	regex := regexp.MustCompile(`.*RESPONSE BODY:\n(.*)`)
	return regex.FindSubmatch([]byte(body))[1]
}
