package fakes

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Tokenizer struct {
	key []byte
}

func NewTokenizer(key string) Tokenizer {
	return Tokenizer{
		key: []byte(key),
	}
}

func (t Tokenizer) Encrypt(token Token) string {
	crypt := jwt.New(jwt.SigningMethodHS256)
	crypt.Claims = token.ToClaims()
	encrypted, err := crypt.SignedString(t.key)
	if err != nil {
		panic(err)
	}

	return encrypted
}

func (t Tokenizer) Decrypt(encryptedToken string) Token {
	token, err := jwt.Parse(encryptedToken, jwt.Keyfunc(func(*jwt.Token) (interface{}, error) {
		return t.key, nil
	}))
	if err != nil {
		panic(err)
	}

	return NewTokenFromClaims(token.Claims)
}

func (t Tokenizer) Validate(token, expected Token) bool {
	if ok := token.HasAudiences(expected.Audiences); !ok {
		return false
	}

	if ok := token.HasScopes(expected.Scopes); !ok {
		return false
	}

	return true
}

type Token struct {
	UserID    string
	ClientID  string
	Scopes    []string
	Audiences []string
}

func NewTokenFromClaims(claims map[string]interface{}) Token {
	token := Token{}

	if userID, ok := claims["user_id"].(string); ok {
		token.UserID = userID
	}

	if clientID, ok := claims["client_id"].(string); ok {
		token.ClientID = clientID
	}

	if scopes, ok := claims["scope"].([]interface{}); ok {
		var s []string
		for _, scope := range scopes {
			s = append(s, scope.(string))
		}

		token.Scopes = s
	}

	if audiences, ok := claims["aud"].(string); ok {
		token.Audiences = strings.Split(audiences, " ")
	}

	return token
}

func (t Token) ToClaims() map[string]interface{} {
	claims := make(map[string]interface{})

	if len(t.UserID) > 0 {
		claims["user_id"] = t.UserID
	}

	if len(t.ClientID) > 0 {
		claims["client_id"] = t.ClientID
	}

	claims["scope"] = t.Scopes
	claims["aud"] = strings.Join(t.Audiences, " ")

	return claims
}

func (t Token) HasScopes(scopes []string) bool {
	for _, scope := range scopes {
		if !contains(t.Scopes, scope) {
			return false
		}
	}
	return true
}

func (t Token) HasAudiences(audiences []string) bool {
	for _, audience := range audiences {
		if !contains(t.Audiences, audience) {
			return false
		}
	}
	return true
}

func contains(collection []string, item string) bool {
	for _, elem := range collection {
		if elem == item {
			return true
		}
	}

	return false
}
