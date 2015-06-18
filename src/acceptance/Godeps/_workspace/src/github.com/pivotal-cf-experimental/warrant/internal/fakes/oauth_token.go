package fakes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

func (s *UAAServer) OAuthToken(w http.ResponseWriter, req *http.Request) {
	_, _, ok := req.BasicAuth()
	if !ok {
		s.Error(w, http.StatusUnauthorized, "An Authentication object was not found in the SecurityContext", "unauthorized")
		return
	}

	err := req.ParseForm()
	if err != nil {
		panic(err)
	}
	clientID := req.Form.Get("client_id")

	scopes := []string{"scim.write","scim.read","password.write"}
	token := s.tokenizer.Encrypt(Token{
		ClientID:  clientID,
		Scopes:    scopes,
		Audiences: []string{"scim","password"},
	})

	response, err := json.Marshal(documents.TokenResponse{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   5000,
		Scope:       strings.Join(scopes, " "),
		JTI:         GenerateID(),
	})

	w.Write(response)
}
