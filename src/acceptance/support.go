package acceptance

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"bitbucket.org/chrj/smtpd"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo/v2"
)

var trace = os.Getenv("TRACE") != ""

type TestSuiteContext struct {
	NotificationsAppGUID        string
	NotificationsOrg            string
	NotificationsSpace          string
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
	TestClientSenderAuthorities []string
	TestClientSenderGrantTypes  []string
	NotificationsDomain         string
	UAADomain                   string
	CCDomain                    string
	LogToken                    string
	Deliveries                  []smtpd.Envelope
}

type GUIDResponse struct {
	Resources []struct {
		MetaData struct {
			GUID string `json:"guid"`
		} `json:"metadata"`
	}
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

	if trace {
		fmt.Println(strings.Join(parts, " "))
	}

	cmd := exec.Command(command, arguments...)
	output, err := cmd.CombinedOutput()
	if trace {
		fmt.Println(string(output))
	}

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

func freePort() string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		Fail(err.Error(), 1)
	}
	defer listener.Close()

	address := listener.Addr().String()
	addressParts := strings.SplitN(address, ":", 2)
	return addressParts[1]
}
