package services

import "github.com/cloudfoundry/notifications-release/src/notifications/v81/v1/models"

type DatabaseInterface interface {
	models.DatabaseInterface
}

type ConnectionInterface interface {
	models.ConnectionInterface
}
