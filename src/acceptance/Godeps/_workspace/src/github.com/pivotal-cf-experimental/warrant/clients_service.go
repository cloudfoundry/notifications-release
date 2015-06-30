package warrant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/network"
)

type Client struct {
	ID                   string
	Scope                []string
	ResourceIDs          []string
	Authorities          []string
	AuthorizedGrantTypes []string
	AccessTokenValidity  time.Duration
}

type ClientsService struct {
	config Config
}

func NewClientsService(config Config) ClientsService {
	return ClientsService{
		config: config,
	}
}

func (cs ClientsService) Create(client Client, secret, token string) error {
	_, err := newNetworkClient(cs.config).MakeRequest(network.Request{
		Method:        "POST",
		Path:          "/oauth/clients",
		Authorization: network.NewTokenAuthorization(token),
		Body:          network.NewJSONRequestBody(client.ToDocument(secret)),
		AcceptableStatusCodes: []int{http.StatusCreated},
	})
	if err != nil {
		return translateError(err)
	}

	return nil
}

func (cs ClientsService) Get(id, token string) (Client, error) {
	resp, err := newNetworkClient(cs.config).MakeRequest(network.Request{
		Method:                "GET",
		Path:                  fmt.Sprintf("/oauth/clients/%s", id),
		Authorization:         network.NewTokenAuthorization(token),
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return Client{}, translateError(err)
	}

	var document documents.ClientResponse
	err = json.Unmarshal(resp.Body, &document)
	if err != nil {
		return Client{}, MalformedResponseError{err}
	}

	return newClientFromDocument(document), nil
}

func (cs ClientsService) Delete(id, token string) error {
	_, err := newNetworkClient(cs.config).MakeRequest(network.Request{
		Method:                "DELETE",
		Path:                  fmt.Sprintf("/oauth/clients/%s", id),
		Authorization:         network.NewTokenAuthorization(token),
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return translateError(err)
	}

	return nil
}

func (cs ClientsService) GetToken(id, secret string) (string, error) {
	resp, err := newNetworkClient(cs.config).MakeRequest(network.Request{
		Method:        "POST",
		Path:          "/oauth/token",
		Authorization: network.NewBasicAuthorization(id, secret),
		Body: network.NewFormRequestBody(url.Values{
			"client_id":  []string{id},
			"grant_type": []string{"client_credentials"},
		}),
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return "", translateError(err)
	}

	var response documents.TokenResponse
	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		return "", MalformedResponseError{err}
	}

	return response.AccessToken, nil
}

func newClientFromDocument(document documents.ClientResponse) Client {
	return Client{
		ID:                   document.ClientID,
		Scope:                document.Scope,
		ResourceIDs:          document.ResourceIDs,
		Authorities:          document.Authorities,
		AuthorizedGrantTypes: document.AuthorizedGrantTypes,
		AccessTokenValidity:  time.Duration(document.AccessTokenValidity) * time.Second,
	}
}

func (c Client) ToDocument(secret string) documents.CreateClientRequest {
	client := documents.CreateClientRequest{
		ClientID:             c.ID,
		ClientSecret:         secret,
		AccessTokenValidity:  int(c.AccessTokenValidity.Seconds()),
		Scope:                make([]string, 0),
		ResourceIDs:          make([]string, 0),
		Authorities:          make([]string, 0),
		AuthorizedGrantTypes: make([]string, 0),
	}
	client.Scope = append(client.Scope, c.Scope...)
	client.ResourceIDs = append(client.ResourceIDs, c.ResourceIDs...)
	client.Authorities = append(client.Authorities, c.Authorities...)
	client.AuthorizedGrantTypes = append(client.AuthorizedGrantTypes, c.AuthorizedGrantTypes...)

	return client
}
