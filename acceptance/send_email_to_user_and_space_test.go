package acceptance_test

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "strings"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

type GUIDResponse struct {
    Resources []struct {
        MetaData struct {
            GUID string `json:guid`
        } `json:metadata`
    }
}

var _ = Describe("SendEmailToUserAndSpace", func() {
    var uaacUserID string
    var orgGUIDResponse GUIDResponse
    var spaceGUIDResponse GUIDResponse
    var uaacAdminClientUserName string
    var uaacAdminClientPassword string

    BeforeSuite(func() {
        // LOAD ENVIRONMENT LOGIN INFO

        uaacAdminClientUserName = os.Getenv("UAAC_ADMIN_CLIENT_USERNAME")
        uaacAdminClientPassword = os.Getenv("UAAC_ADMIN_CLIENT_PASSWORD")
        cfAdminUsername := os.Getenv("CF_ADMIN_USERNAME")
        cfAdminPassword := os.Getenv("CF_ADMIN_PASSWORD")

        // LOGIN AS A CF USER
        command := exec.Command("cf", "login", "-u", cfAdminUsername, "-p", cfAdminPassword)
        err := command.Run()
        if err != nil {
            panic("Couldn't login to CF, check your credentials")
        }

        // CREATE A USER AND GRAB ITS TOKEN
        command = exec.Command("cf", "create-user", "notificationsTestUser", "password")
        command.Run()
        command = exec.Command("cf", "create-org", "notificationsTestOrg")
        command.Run()
        command = exec.Command("cf", "create-space", "notificationsTestSpace", "-o", "notificationsTestOrg")
        command.Run()

        command = exec.Command("uaac", "token", "client", "get", uaacAdminClientUserName, "-s", uaacAdminClientPassword)
        err = command.Run()
        if err != nil {
            panic("Coudln't login to UAAC, check your credentials")
        }

        command = exec.Command("uaac", "user", "get", "notificationsTestUser", "-a", "id")
        data, _ := command.Output()
        uaacUserID = strings.Split(string(data), ":")[1]
        uaacUserID = strings.TrimSpace(uaacUserID)

        // GET A CLIENT WITH THE RIGHT SCOPES
        command = exec.Command("uaac", "client", "add", "notifications-sender", "--scope", "notifications.write", "--authorities",
            "notifications.write", "-s", "secret", "--authorized_grant_types", "client_credentials")
        command.Run()
        command = exec.Command("uaac", "token", "client", "get", "notifications-sender", "-s", "secret")
        command.Run()
    })

    AfterSuite(func() {
        command := exec.Command("cf", "delete-user", "notificationsTestUser", "-f")
        command.Run()
        command = exec.Command("cf", "delete-space", "notificationsTestSpace", "-f")
        command.Run()
        command = exec.Command("cf", "delete-org", "notificationsTestOrg", "-f")
        command.Run()

        command = exec.Command("uaac", "token", "client", "get", uaacAdminClientUserName, "-s", uaacAdminClientPassword)
        command.Run()

        command = exec.Command("uaac", "client", "delete", "notifications-sender")
        command.Run()
    })

    Describe("sending an email to a user", func() {
        It("returns a 200", func() {
            // SEND A NOTIFICATION TO A USER
            notificationToUserURL := fmt.Sprintf("%s/users/%s", os.Getenv("NOTIFICATIONS_DOMAIN"), uaacUserID)
            command := exec.Command("uaac", "curl", notificationToUserURL, "-X", "POST", "--data", `{"kind_id":"test_notification", "text":"this is a test"}`)
            data, _ := command.Output()

            // VERIFY 200 RESPONSE
            Expect(data).To(ContainSubstring("200 OK"))
        })
    })

    Describe("sending an email to a space", func() {
        It("returns a 200", func() {
            // PUT A USER IN A SPACE
            command := exec.Command("cf", "set-space-role", "notificationsTestUser", "notificationsTestOrg", "notificationsTestSpace", "SpaceDeveloper")
            command.Run()
            command = exec.Command("cf", "login", "-u", "notificationsTestUser", "-p", "password")
            command.Run()

            // GRAB IDs FOR CURL REQUESTS
            command = exec.Command("cf", "curl", "/v2/organizations")
            data, _ := command.Output()
            err := json.Unmarshal(data, &orgGUIDResponse)
            if err != nil {
                panic(err)
            }

            orgGUID := orgGUIDResponse.Resources[0].MetaData.GUID

            command = exec.Command("cf", "curl", fmt.Sprintf("/v2/organizations/%s/spaces", orgGUID))
            data, _ = command.Output()
            err = json.Unmarshal(data, &spaceGUIDResponse)
            if err != nil {
                panic(err)
            }
            spaceGUID := spaceGUIDResponse.Resources[0].MetaData.GUID

            // SEND A NOTIFICATION TO A SPACE
            notificationToSpaceURL := fmt.Sprintf("%s/spaces/%s", os.Getenv("NOTIFICATIONS_DOMAIN"), spaceGUID)
            command = exec.Command("uaac", "curl", notificationToSpaceURL, "-X", "POST", "--data", `{"kind_id":"test_notification", "text":"this is a test"}`)
            data, _ = command.Output()

            // VERIFY 200 RESPONSE
            Expect(data).To(ContainSubstring("200 OK"))
        })
    })

})
