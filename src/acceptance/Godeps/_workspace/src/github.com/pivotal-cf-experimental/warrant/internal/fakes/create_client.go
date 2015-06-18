package fakes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

func (s *UAAServer) CreateClient(w http.ResponseWriter, req *http.Request) {
	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if ok := s.ValidateToken(token, []string{"clients"}, []string{"clients.write"}); !ok {
		s.Error(w, http.StatusUnauthorized, "Full authentication is required to access this resource", "unauthorized")
		return
	}

	var document documents.CreateClientRequest
	err := json.NewDecoder(req.Body).Decode(&document)
	if err != nil {
		panic(err)
	}

	client := newClientFromDocument(document)
	if err := client.Validate(); err != nil {
		s.Error(w, http.StatusBadRequest, err.Error(), "invalid_client")
		return
	}

	s.clients.Add(client)

	response, err := json.Marshal(client.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
