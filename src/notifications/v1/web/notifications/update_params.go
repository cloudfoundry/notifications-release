package notifications

import (
	"io"

	"github.com/cloudfoundry/notifications-release/src/notifications/v81/v1/models"
	"github.com/cloudfoundry/notifications-release/src/notifications/v81/v1/web/webutil"
	"github.com/cloudfoundry/notifications-release/src/notifications/v81/valiant"
)

type NotificationUpdateParams struct {
	Description string `json:"description" validate-required:"true"`
	Critical    bool   `json:"critical"    validate-required:"true"`
	TemplateID  string `json:"template"    validate-required:"true"`
}

func NewNotificationParams(body io.Reader) (NotificationUpdateParams, error) {
	var params NotificationUpdateParams

	validator := valiant.NewValidator(body)
	err := validator.Validate(&params)
	if err != nil {
		switch err.(type) {
		case valiant.RequiredFieldError:
			return params, webutil.ValidationError{Err: err}
		default:
			return params, webutil.ParseError{}
		}
	}
	return params, nil
}

func (params NotificationUpdateParams) ToModel(clientID, notificationID string) models.Kind {
	return models.Kind{
		Description: params.Description,
		Critical:    params.Critical,
		TemplateID:  params.TemplateID,
		ClientID:    clientID,
		ID:          notificationID,
	}
}
