package documents

import (
	"encoding/json"
	"time"
)

type CreateUserRequest struct {
	UserName string   `json:"userName"`
	Name     UserName `json:"name"`
	Emails   []Email  `json:"emails"`
}

type UpdateUserRequest struct {
	Schemas    []string `json:"schemas"`
	ID         string   `json:"id"`
	UserName   string   `json:"userName"`
	ExternalID string   `json:"externalId"`
	Name       UserName `json:"name"`
	Emails     []Email  `json:"emails"`
	Meta       Meta     `json:"meta"`
}

type UserResponse struct {
	Schemas    []string `json:"schemas"`
	ID         string   `json:"id"`
	ExternalID string   `json:"externalId"`
	UserName   string   `json:"userName"`
	Name       UserName `json:"name"`
	Emails     []Email  `json:"emails"`
	Meta       Meta     `json:"meta"`
	Groups     []Group  `json:"groups"`
	Active     bool     `json:"active"`
	Verified   bool     `json:"verified"`
	Origin     string   `json:"origin"`
}

type UserListResponse struct {
	Resources    []UserResponse `json:"resources"`
	StartIndex   int            `json:"startIndex"`
	ItemsPerPage int            `json:"itemsPerPage"`
	TotalResults int            `json:"totalResults"`
	Schemas      []string       `json:"schemas"`
}

type UserName struct {
	Formatted  string `json:"formatted"`
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
	MiddleName string `json:"middleName"`
}

type Email struct {
	Value string `json:"value"`
}

type Meta struct {
	Version      int       `json:"version"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

// TODO: UAA team is investigating this hack as a possible bug
func (m Meta) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"version":      m.Version,
		"created":      m.Created.Format("2006-01-02T15:04:05.000Z"),
		"lastModified": m.LastModified.Format("2006-01-02T15:04:05.000Z"),
	})
}

type Group struct{}
