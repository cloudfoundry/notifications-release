package fakes

import (
	"encoding/json"
	"fmt"
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

	filter := query.Get("filter")
	matches := regexp.MustCompile(`(.*) (.*) '(.*)'$`).FindStringSubmatch(filter)
	parameter := matches[1]
	operator := matches[2]
	value := matches[3]

	if !validParameter(parameter) {
		s.Error(w, http.StatusBadRequest, fmt.Sprintf("Invalid filter expression: [%s]", filter), "scim")
		return
	}

	if !validOperator(operator) {
		s.Error(w, http.StatusBadRequest, fmt.Sprintf("Invalid filter expression: [%s]", filter), "scim")
		return
	}

	user, _ := s.users.Get(value)

	list := UsersList{user}

	response, err := json.Marshal(list.ToDocument())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func validParameter(parameter string) bool {
	for _, p := range []string{"id"} {
		if parameter == p {
			return true
		}
	}

	return false
}

func validOperator(operator string) bool {
	for _, o := range []string{"eq"} {
		if operator == o {
			return true
		}
	}

	return false
}
