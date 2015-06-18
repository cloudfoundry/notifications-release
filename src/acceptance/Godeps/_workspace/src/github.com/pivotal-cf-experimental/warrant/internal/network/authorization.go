package network

import (
	"encoding/base64"
	"fmt"
)

type authorization interface {
	Authorization() string
}

func NewTokenAuthorization(token string) tokenAuthorization {
	return tokenAuthorization(token)
}

type tokenAuthorization string

func (a tokenAuthorization) Authorization() string {
	return fmt.Sprintf("Bearer %s", a)
}

func NewBasicAuthorization(username, password string) basicAuthorization {
	return basicAuthorization{
		Username: username,
		Password: password,
	}
}

type basicAuthorization struct {
	Username string
	Password string
}

func (b basicAuthorization) Authorization() string {
	auth := b.Username + ":" + b.Password
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
}
