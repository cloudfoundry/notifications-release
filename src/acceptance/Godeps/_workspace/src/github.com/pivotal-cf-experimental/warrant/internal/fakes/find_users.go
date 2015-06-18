package fakes

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

type UsersList []User

func (ul UsersList) ToDocument() documents.UserListResponse {
	doc := documents.UserListResponse{
		ItemsPerPage: 100,
		StartIndex:   1,
		TotalResults: len(ul),
		Schemas:      Schemas,
	}

	for _, user := range ul {
		doc.Resources = append(doc.Resources, user.ToDocument())
	}

	return doc
}

func (s *UAAServer) FindUsers(w http.ResponseWriter, req *http.Request) {
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		panic(err)
	}

	matches := regexp.MustCompile(`id eq '(.*)'$`).FindStringSubmatch(query.Get("filter"))
	id := matches[1]
	user, _ := s.users.Get(id)

	list := UsersList{user}

	response, err := json.Marshal(list.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
