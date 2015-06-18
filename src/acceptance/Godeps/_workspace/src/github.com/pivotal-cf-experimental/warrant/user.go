package warrant

import (
	"time"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

type User struct {
	ID            string
	ExternalID    string
	UserName      string
	FormattedName string
	FamilyName    string
	GivenName     string
	MiddleName    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Version       int
	Emails        []string
	Groups        []Group
	Active        bool
	Verified      bool
	Origin        string
}

func newUserFromResponse(config Config, response documents.UserResponse) User {
	var emails []string
	for _, email := range response.Emails {
		emails = append(emails, email.Value)
	}

	return User{
		ID:            response.ID,
		ExternalID:    response.ExternalID,
		UserName:      response.UserName,
		FormattedName: response.Name.Formatted,
		FamilyName:    response.Name.FamilyName,
		GivenName:     response.Name.GivenName,
		MiddleName:    response.Name.MiddleName,
		Emails:        emails,
		CreatedAt:     response.Meta.Created,
		UpdatedAt:     response.Meta.LastModified,
		Active:        response.Active,
		Verified:      response.Verified,
		Origin:        response.Origin,
	}
}
