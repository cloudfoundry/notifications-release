package warrant

import (
	"io"

	"github.com/pivotal-cf-experimental/warrant/internal/network"
)

type Config struct {
	Host          string
	SkipVerifySSL bool
	TraceWriter   io.Writer
}

type Warrant struct {
	config  Config
	Users   UsersService
	Clients ClientsService
	Tokens  TokensService
}

func New(config Config) Warrant {
	return Warrant{
		config:  config,
		Users:   NewUsersService(config),
		Clients: NewClientsService(config),
		Tokens:  NewTokensService(config),
	}
}

func newNetworkClient(config Config) network.Client {
	return network.NewClient(network.Config{
		Host:          config.Host,
		SkipVerifySSL: config.SkipVerifySSL,
		TraceWriter:   config.TraceWriter,
	})
}
