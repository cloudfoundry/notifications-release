package fakes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

func (s *UAAServer) UpdateUser(w http.ResponseWriter, req *http.Request) {
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

	var document documents.UpdateUserRequest
	err = json.Unmarshal(requestBody, &document)
	if err != nil {
		panic(err)
	}

	user := newUserFromUpdateDocument(document)

	matches := regexp.MustCompile(`/Users/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	existingUser, ok := s.users.Get(id)
	if !ok {
		s.Error(w, http.StatusNotFound, fmt.Sprintf("User %s does not exist", user.ID), "scim_resource_not_found")
		return
	}

	version, err := strconv.ParseInt(req.Header.Get("If-Match"), 10, 64)
	if err != nil || existingUser.Version != int(version) {
		s.Error(w, http.StatusBadRequest, "Missing If-Match for PUT", "scim")
		return
	}

	s.users.Update(user)

	response, err := json.Marshal(user.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
