package fakes

import (
	"errors"
	"fmt"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

var ValidGrantTypes = []string{"implicit", "refresh_token", "authorization_code", "client_credentials", "password"}

type Client struct {
	ID                   string
	Secret               string
	Scope                []string
	ResourceIDs          []string
	Authorities          []string
	AuthorizedGrantTypes []string
	AccessTokenValidity  int
}

func newClientFromDocument(document documents.CreateClientRequest) Client {
	return Client{
		ID:                   document.ClientID,
		Secret:               document.ClientSecret,
		Scope:                document.Scope,
		ResourceIDs:          document.ResourceIDs,
		Authorities:          document.Authorities,
		AuthorizedGrantTypes: document.AuthorizedGrantTypes,
		AccessTokenValidity:  document.AccessTokenValidity,
	}
}

func (c Client) ToDocument() documents.ClientResponse {
	return documents.ClientResponse{
		ClientID:             c.ID,
		Scope:                c.Scope,
		ResourceIDs:          c.ResourceIDs,
		Authorities:          c.Authorities,
		AuthorizedGrantTypes: c.AuthorizedGrantTypes,
		AccessTokenValidity:  c.AccessTokenValidity,
	}
}

func (c Client) Validate() error {
	for _, grantType := range c.AuthorizedGrantTypes {
		if !contains(ValidGrantTypes, grantType) {
			msg := fmt.Sprintf("%s is not an allowed grant type. Must be one of: %v", grantType, ValidGrantTypes)
			return errors.New(msg)
		}
	}

	return nil
}
