package main

import (
	"log"

	"github.com/cloudfoundry/notifications-release/src/notifications/v81/application"
)

func main() {
	env, err := application.NewEnvironment()
	if err != nil {
		log.Fatalf("CRASHING: %s\n", err)
	}

	dbp := application.NewDBProvider(env)
	app := application.New(env, dbp)
	defer app.Crash()

	app.Run()
}
