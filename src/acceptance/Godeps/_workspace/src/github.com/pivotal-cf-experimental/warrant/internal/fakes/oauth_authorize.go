package fakes

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (s *UAAServer) OAuthAuthorize(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Accept") != "application/json" {
		w.Header().Set("Location", fmt.Sprintf("%s/login", s.URL()))
		w.WriteHeader(http.StatusFound)
		return
	}

	requestQuery := req.URL.Query()
	clientID := requestQuery.Get("client_id")
	responseType := requestQuery.Get("response_type")

	if clientID != "cf" {
		w.Header().Set("Location", fmt.Sprintf("%s/login", s.URL()))
		w.WriteHeader(http.StatusFound)
		return
	}

	if responseType != "token" {
		w.Header().Set("Location", fmt.Sprintf("%s/login", s.URL()))
		w.WriteHeader(http.StatusFound)
		return
	}

	req.ParseForm()
	userName := req.Form.Get("username")

	user, ok := s.users.GetByName(userName)
	if !ok {
		s.Error(w, http.StatusNotFound, fmt.Sprintf("User %s does not exist", userName), "scim_resource_not_found")
		return
	}

	if req.Form.Get("source") != "credentials" {
		w.Header().Set("Location", fmt.Sprintf("%s/login", s.URL()))
		w.WriteHeader(http.StatusFound)
		return
	}

	if req.Form.Get("password") != user.Password {
		w.Header().Set("Location", fmt.Sprintf("%s/login", s.URL()))
		w.WriteHeader(http.StatusFound)
		return
	}

	scopes := strings.Join(s.defaultScopes, " ")

	token := s.tokenizer.Encrypt(Token{
		UserID:    user.ID,
		Scopes:    s.defaultScopes,
		Audiences: []string{},
	})

	redirectURI := requestQuery.Get("redirect_uri")

	query := url.Values{
		"token_type":   []string{"bearer"},
		"access_token": []string{token},
		"expires_in":   []string{"599"},
		"scope":        []string{scopes},
		"jti":          []string{"ad0efc96-ed29-43ef-be75-85a4b4f105b5"},
	}
	location := fmt.Sprintf("%s#%s", redirectURI, query.Encode())

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusFound)
}
