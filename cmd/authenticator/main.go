package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"go.uber.org/zap"

	"github.com/migotom/cell-centre-services/db"
	authDelivery "github.com/migotom/cell-centre-services/pkg/components/auth/delivery/grpc"
	employeeRepository "github.com/migotom/cell-centre-services/pkg/components/employee/repository"
	"github.com/migotom/cell-centre-services/pkg/services/authenticator"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var config authenticator.Config

	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend("/etc/cell-centre/authenticator/config.toml"),
	)
	if err := loader.Load(context.Background(), &config); err != nil {
		log.Fatal("Error during config reading", zap.Error(err))
	}

	dbClient, db, err := db.ConnectMongoDB(context.Background(), config.DatabaseAddress, config.DatabaseName)
	if err != nil {
		log.Fatal("Can't connect to database", zap.Error(err))
	}
	defer func() {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			log.Fatal("Can't safely disconnect from database", zap.Error(err))
		}
	}()

	authenticator := authenticator.NewAuthenticator(
		log,
		&config,
		authDelivery.NewAuthenticateDelivery(
			log,
			employeeRepository.NewEmployeeRepository(db),
		),
	)
	authenticator.Listen()
}
