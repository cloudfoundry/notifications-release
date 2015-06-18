package warrant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/network"
)

// TODO: Score password strength
// TODO: Verify a user
// TODO: Query for user info
// TODO: Convert user ids to names

const Schema = "urn:scim:schemas:core:1.0"

var Schemas = []string{Schema}

type UsersService struct {
	config Config
}

func NewUsersService(config Config) UsersService {
	return UsersService{
		config: config,
	}
}

func (us UsersService) Create(username, email, token string) (User, error) {
	resp, err := newNetworkClient(us.config).MakeRequest(network.Request{
		Method:        "POST",
		Path:          "/Users",
		Authorization: network.NewTokenAuthorization(token),
		Body: network.NewJSONRequestBody(documents.CreateUserRequest{
			UserName: username,
			Emails: []documents.Email{
				{Value: email},
			},
		}),
		AcceptableStatusCodes: []int{http.StatusCreated},
	})
	if err != nil {
		return User{}, translateError(err)
	}

	var response documents.UserResponse
	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		panic(err)
	}

	return newUserFromResponse(us.config, response), nil
}

func (us UsersService) Get(id, token string) (User, error) {
	resp, err := newNetworkClient(us.config).MakeRequest(network.Request{
		Method:                "GET",
		Path:                  fmt.Sprintf("/Users/%s", id),
		Authorization:         network.NewTokenAuthorization(token),
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return User{}, translateError(err)
	}

	var response documents.UserResponse
	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		panic(err)
	}

	return newUserFromResponse(us.config, response), nil
}

func (us UsersService) Delete(id, token string) error {
	_, err := newNetworkClient(us.config).MakeRequest(network.Request{
		Method:                "DELETE",
		Path:                  fmt.Sprintf("/Users/%s", id),
		Authorization:         network.NewTokenAuthorization(token),
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return translateError(err)
	}

	return nil
}

func (us UsersService) Update(user User, token string) (User, error) {
	resp, err := newNetworkClient(us.config).MakeRequest(network.Request{
		Method:        "PUT",
		Path:          fmt.Sprintf("/Users/%s", user.ID),
		Authorization: network.NewTokenAuthorization(token),
		IfMatch:       strconv.Itoa(user.Version),
		Body:          network.NewJSONRequestBody(newUpdateUserDocumentFromUser(user)),
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return User{}, translateError(err)
	}

	var response documents.UserResponse
	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		panic(err)
	}

	return newUserFromResponse(us.config, response), nil
}

func (us UsersService) SetPassword(id, password, token string) error {
	_, err := newNetworkClient(us.config).MakeRequest(network.Request{
		Method:        "PUT",
		Path:          fmt.Sprintf("/Users/%s/password", id),
		Authorization: network.NewTokenAuthorization(token),
		Body: network.NewJSONRequestBody(documents.SetPasswordRequest{
			Password: password,
		}),
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return translateError(err)
	}

	return nil
}

func (us UsersService) ChangePassword(id, oldPassword, password, token string) error {
	_, err := newNetworkClient(us.config).MakeRequest(network.Request{
		Method:        "PUT",
		Path:          fmt.Sprintf("/Users/%s/password", id),
		Authorization: network.NewTokenAuthorization(token),
		Body: network.NewJSONRequestBody(documents.ChangePasswordRequest{
			OldPassword: oldPassword,
			Password:    password,
		}),
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return translateError(err)
	}

	return nil
}

func (us UsersService) GetToken(username, password string) (string, error) {
	query := url.Values{
		"client_id":     []string{"cf"},
		"redirect_uri":  []string{"https://uaa.cloudfoundry.com/redirect/cf"},
		"response_type": []string{"token"},
	}

	requestPath := url.URL{
		Path:     "/oauth/authorize",
		RawQuery: query.Encode(),
	}
	req := network.Request{
		Method: "POST",
		Path:   requestPath.String(),
		Body: network.NewFormRequestBody(url.Values{
			"username": []string{username},
			"password": []string{password},
			"source":   []string{"credentials"},
		}),
		AcceptableStatusCodes: []int{http.StatusFound},
		DoNotFollowRedirects:  true,
	}

	resp, err := newNetworkClient(us.config).MakeRequest(req)
	if err != nil {
		return "", translateError(err)
	}

	locationURL, err := url.Parse(resp.Headers.Get("Location"))
	if err != nil {
		return "", err
	}

	locationQuery, err := url.ParseQuery(locationURL.Fragment)
	if err != nil {
		return "", err
	}

	return locationQuery.Get("access_token"), nil
}

type UsersQuery struct {
	Filter string
}

func (us UsersService) Find(query UsersQuery, token string) ([]User, error) {
	requestPath := url.URL{
		Path: "/Users",
		RawQuery: url.Values{
			"filter": []string{query.Filter},
		}.Encode(),
	}

	resp, err := newNetworkClient(us.config).MakeRequest(network.Request{
		Method:                "GET",
		Path:                  requestPath.String(),
		Authorization:         network.NewTokenAuthorization(token),
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		panic(err)
	}

	var response documents.UserListResponse
	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		panic(err)
	}

	var userList []User
	for _, userResponse := range response.Resources {
		userList = append(userList, newUserFromResponse(us.config, userResponse))
	}

	return userList, err
}

func newUpdateUserDocumentFromUser(user User) documents.UpdateUserRequest {
	var emails []documents.Email
	for _, email := range user.Emails {
		emails = append(emails, documents.Email{
			Value: email,
		})
	}

	return documents.UpdateUserRequest{
		Schemas:    Schemas,
		ID:         user.ID,
		UserName:   user.UserName,
		ExternalID: user.ExternalID,
		Name: documents.UserName{
			Formatted:  user.FormattedName,
			FamilyName: user.FamilyName,
			GivenName:  user.GivenName,
			MiddleName: user.MiddleName,
		},
		Emails: emails,
		Meta: documents.Meta{
			Version:      user.Version,
			Created:      user.CreatedAt,
			LastModified: user.UpdatedAt,
		},
	}
}
