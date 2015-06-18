package fakes

import (
	"net/http"
	"regexp"
	"strings"
)

func (s *UAAServer) DeleteUser(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := s.ValidateToken(token, []string{"scim"}, []string{"scim.write"}); !ok {
		s.Error(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	matches := regexp.MustCompile(`/Users/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	if ok := s.users.Delete(id); !ok {
		s.Error(w, http.StatusNotFound, "User non-existant-user-guid does not exist", "scim_resource_not_found")
		return
	}

	w.WriteHeader(http.StatusOK)
}
