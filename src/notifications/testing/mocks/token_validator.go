package mocks

import (
	gouaa "github.com/cloudfoundry-community/go-uaa"
	"github.com/golang-jwt/jwt/v5"
)

type TokenValidator struct {
	ParseCall struct {
		Receives struct {
			Token string
		}

		Returns struct {
			Token *jwt.Token
			Error error
		}
	}
}

func (t *TokenValidator) Parse(token string) (*jwt.Token, error) {
	t.ParseCall.Receives.Token = token
	return t.ParseCall.Returns.Token, t.ParseCall.Returns.Error
}

type KeyFetcher struct {
	GetSigningKeysCall struct {
		Called  bool
		Returns struct {
			Keys  []gouaa.JWK
			Error error
		}
	}
}

func (f *KeyFetcher) TokenKeys() ([]gouaa.JWK, error) {
	f.GetSigningKeysCall.Called = true
	return f.GetSigningKeysCall.Returns.Keys, f.GetSigningKeysCall.Returns.Error
}
