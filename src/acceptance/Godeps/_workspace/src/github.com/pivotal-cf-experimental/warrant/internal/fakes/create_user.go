package fakes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

func (s *UAAServer) CreateUser(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := s.ValidateToken(token, []string{"scim"}, []string{"scim.write"}); !ok {
		s.Error(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		if contentType == "" {
			contentType = http.DetectContentType(requestBody)
		}
		s.Error(w, http.StatusBadRequest, fmt.Sprintf("Content type '%s' not supported", contentType), "scim")
		return
	}

	var document documents.CreateUserRequest
	err = json.Unmarshal(requestBody, &document)
	if err != nil {
		panic(err)
	}

	user := newUserFromCreateDocument(document)
	if err := user.Validate(); err != nil {
		s.Error(w, http.StatusBadRequest, err.Error(), "invalid_scim_resource")
		return
	}
	s.users.Add(user)

	response, err := json.Marshal(user.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
