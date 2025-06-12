package notify

import "github.com/cloudfoundry/notifications-release/src/notifications/v81/v1/services"

type DatabaseInterface interface {
	services.DatabaseInterface
}

type ConnectionInterface interface {
	services.ConnectionInterface
}
