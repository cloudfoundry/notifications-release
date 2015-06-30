package warrant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/network"
)

type Token struct {
	ClientID string   `json:"client_id"`
	UserID   string   `json:"user_id"`
	Scopes   []string `json:"scope"`
}

type TokensService struct {
	config Config
}

type SigningKey struct {
	Algorithm string
	Value     string
}

func NewTokensService(config Config) TokensService {
	return TokensService{
		config: config,
	}
}

func (ts TokensService) Decode(token string) (Token, error) {
	segments := strings.Split(token, ".")
	if len(segments) != 3 {
		return Token{}, InvalidTokenError{fmt.Errorf("invalid number of segments in token (%d/3)", len(segments))}
	}

	claims, err := jwt.DecodeSegment(segments[1])
	if err != nil {
		return Token{}, InvalidTokenError{fmt.Errorf("claims cannot be decoded: %s", err)}
	}

	t := Token{}
	err = json.Unmarshal(claims, &t)
	if err != nil {
		return Token{}, InvalidTokenError{fmt.Errorf("token cannot be parsed: %s", err)}
	}

	return t, nil
}

func (ts TokensService) GetSigningKey() (SigningKey, error) {
	resp, err := newNetworkClient(ts.config).MakeRequest(network.Request{
		Method: "GET",
		Path:   "/token_key",
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return SigningKey{}, translateError(err)
	}

	var response documents.TokenKeyResponse
	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		return SigningKey{}, MalformedResponseError{err}
	}

	return SigningKey{Algorithm: response.Alg, Value: response.Value}, nil
}
