package warrant

import (
	"encoding/json"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
	Scopes 	 []string `json:"scope"`
}

type TokensService struct{}

func NewTokensService(config Config) TokensService {
	return TokensService{}
}

func (ts TokensService) Decode(token string) (Token, error) {
	segments := strings.Split(token, ".")
	claims, err := jwt.DecodeSegment(segments[1])
	if err != nil {
		panic(err)
	}

	t := Token{}
	err = json.Unmarshal(claims, &t)
	if err != nil {
		panic(err)
	}

	return t, nil
}
