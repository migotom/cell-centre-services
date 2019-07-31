package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	stan "github.com/nats-io/go-nats-streaming"
	"go.uber.org/zap"

	"github.com/migotom/cell-centre-services/db"
	authDelivery "github.com/migotom/cell-centre-services/pkg/components/auth/delivery/grpc"
	employeeRepository "github.com/migotom/cell-centre-services/pkg/components/employee/repository"
	"github.com/migotom/cell-centre-services/pkg/components/event"
	roleRepository "github.com/migotom/cell-centre-services/pkg/components/role/repository"
	"github.com/migotom/cell-centre-services/pkg/services/eventstore"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var config eventstore.Config

	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend("/etc/cell-centre/eventstore/config.toml"),
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

	eventsStreaming := event.NewEventsStreaming()
	err = eventsStreaming.Connect(
		config.NATSClusterID,
		stan.NatsURL(config.NATSURL))
	if err != nil {
		log.Fatal("Failed to connect to NATS service", zap.Error(err))
	}

	employeeRepository := employeeRepository.NewEmployeeRepository(db)
	roleRepository := roleRepository.NewRoleRepository(db)
	authDelivery := authDelivery.NewAuthenticateDelivery(log, employeeRepository)

	eventStore := eventstore.NewEventStore(
		log,
		&config,
		eventsStreaming,
		authDelivery,
		employeeRepository,
		roleRepository,
	)
	eventStore.Listen()
}
