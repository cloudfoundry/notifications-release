package helpers

import (
	"github.com/cloudfoundry/notifications-release/src/notifications/v81/application"
	"github.com/cloudfoundry/notifications-release/src/notifications/v81/db"
	"github.com/cloudfoundry/notifications-release/src/notifications/v81/v1/models"
)

func TruncateTables(database *db.DB) {
	env, err := application.NewEnvironment()
	if err != nil {
		panic(err)
	}

	dbMigrator := models.DatabaseMigrator{}
	dbMigrator.Migrate(database.RawConnection(), env.ModelMigrationsPath)
	models.Setup(database)

	connection := database.Connection().(*db.Connection)
	err = connection.TruncateTables()
	if err != nil {
		panic(err)
	}
}
