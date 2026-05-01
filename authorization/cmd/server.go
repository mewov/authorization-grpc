package main

import (
	"log/slog"

	"github.com/gox7/authorizathion/internal/repository"
	"github.com/gox7/authorizathion/internal/services"
	"github.com/gox7/authorizathion/internal/transport"
	"github.com/gox7/authorizathion/models"
	"google.golang.org/grpc"
)

func main() {
	config := new(models.LocalConfig)
	services.NewConfig(config)

	serverLog := new(slog.Logger)
	databaseLog := new(slog.Logger)
	services.NewLoggger("server", &serverLog)
	services.NewLoggger("database", &databaseLog)

	postgres := new(repository.Postgres)
	repository.NewPostgres(config, databaseLog, postgres)
	postgres.Migration()

	sessions := new(services.SessionsService)
	authorizathion := new(services.AuthorizathionService)
	services.NewSessions(postgres, sessions)
	services.NewAuthorizathion(postgres, authorizathion)

	server := grpc.NewServer()
	transport.SetResource(config, authorizathion, sessions)
	transport.Register(serverLog, server)
	go transport.Listen(serverLog, server)
	transport.GracefulShutdown(serverLog, server)
}
