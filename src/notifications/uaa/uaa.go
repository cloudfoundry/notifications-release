package uaa

import (
	"context"
	"fmt"
	"net/url"

	gouaa "github.com/cloudfoundry-community/go-uaa"
	"github.com/golang-jwt/jwt/v5"
	uaaSSOGolang "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type ZonedUAAClient struct {
	clientID       string
	clientSecret   string
	verifySSL      bool
	tokenValidator *TokenValidator
}

func NewZonedUAAClient(clientID, clientSecret string, verifySSL bool, validator *TokenValidator) (client ZonedUAAClient) {
	return ZonedUAAClient{
		clientID:       clientID,
		clientSecret:   clientSecret,
		verifySSL:      verifySSL,
		tokenValidator: validator,
	}
}

func (z ZonedUAAClient) GetTokenKey(uaaHost string) (string, error) {
	uaaClient, err := gouaa.New(
		uaaHost,
		gouaa.WithNoAuthentication(),
		gouaa.WithSkipSSLValidation(!z.verifySSL),
	)
	if err != nil {
		return "", err
	}

	signingKey, err := uaaClient.TokenKey()
	if err != nil {
		return "", err
	}

	return signingKey.Value, nil
}

func (z ZonedUAAClient) GetClientToken(host string) (string, error) {
	uaaClient, err := gouaa.New(
		host,
		gouaa.WithClientCredentials(z.clientID, z.clientSecret, gouaa.JSONWebToken),
		gouaa.WithSkipSSLValidation(!z.verifySSL),
	)
	if err != nil {
		return "", err
	}
	token, err := uaaClient.Token(context.Background())
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func (z ZonedUAAClient) UsersEmailsByIDs(token string, ids ...string) ([]User, error) {
	uaaHost, err := z.tokenHost(token)
	if err != nil {
		return nil, err
	}

	uaaClient := uaaSSOGolang.NewUAA("", uaaHost, z.clientID, z.clientSecret, "")
	uaaClient.VerifySSL = z.verifySSL
	uaaClient.SetToken(token)

	var myUsers []User
	users, err := uaaClient.UsersEmailsByIDs(ids...)
	if err != nil {
		return myUsers, err
	}

	for _, user := range users {
		myUsers = append(myUsers, newUserFromSSOGolangUser(user))
	}

	return myUsers, nil
}

func (z ZonedUAAClient) tokenHost(token string) (string, error) {
	parsedToken, err := z.tokenValidator.Parse(token)
	claims := parsedToken.Claims.(jwt.MapClaims)

	if err != nil {
		return "", err
	}

	tokenIssuerURL, err := url.Parse(claims["iss"].(string))
	if err != nil {
		return "", err
	}

	return tokenIssuerURL.Scheme + "://" + tokenIssuerURL.Host, nil
}

func (z ZonedUAAClient) AllUsers(token string) ([]User, error) {
	uaaHost, err := z.tokenHost(token)
	if err != nil {
		return nil, err
	}

	uaaSSOGolangClient := uaaSSOGolang.NewUAA("", uaaHost, z.clientID, z.clientSecret, "")
	uaaSSOGolangClient.VerifySSL = z.verifySSL
	users, err := uaaSSOGolangClient.AllUsers()

	var myUsers []User
	for _, user := range users {
		myUsers = append(myUsers, newUserFromSSOGolangUser(user))
	}

	return myUsers, err
}

func (z ZonedUAAClient) UsersGUIDsByScope(token string, scope string) ([]string, error) {
	uaaHost, err := z.tokenHost(token)
	if err != nil {
		return nil, err
	}

	uaaSSOGolangClient := uaaSSOGolang.NewUAA("", uaaHost, z.clientID, z.clientSecret, "")
	uaaSSOGolangClient.VerifySSL = z.verifySSL

	return uaaSSOGolangClient.UsersGUIDsByScope(scope)
}

func newUserFromSSOGolangUser(uaaUser uaaSSOGolang.User) User {
	user := User{}
	user.ID = uaaUser.ID
	user.Emails = uaaUser.Emails

	return user
}

type User struct {
	ID     string
	Emails []string
}

type Failure struct {
	code    int
	message string
}

func NewFailure(code int, message []byte) Failure {
	return Failure{
		code:    code,
		message: string(message),
	}
}

func (failure Failure) Code() int {
	return failure.code
}

func (failure Failure) Message() string {
	return failure.message
}

func (failure Failure) Error() string {
	return fmt.Sprintf("UAA Wrapper Failure: %d %s", failure.code, failure.message)
}
