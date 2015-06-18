package fakes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func (s *UAAServer) GetClient(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := s.ValidateToken(token, []string{"clients"}, []string{"clients.read"}); !ok {
		s.Error(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	matches := regexp.MustCompile(`/oauth/clients/(.*)$`).FindStringSubmatch(req.URL.Path)
	id := matches[1]

	client, ok := s.clients.Get(id)
	if !ok {
		s.NotFound(w, fmt.Sprintf("Client %s does not exist", id))
		return
	}

	document := client.ToDocument()
	response, err := json.Marshal(document)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
