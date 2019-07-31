package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	stan "github.com/nats-io/go-nats-streaming"
	"go.uber.org/zap"

	"github.com/migotom/cell-centre-services/db"
	"github.com/migotom/cell-centre-services/pkg/components/event"
	"github.com/migotom/cell-centre-services/pkg/components/event/repository"
	"github.com/migotom/cell-centre-services/pkg/services/eventlogger"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var config eventlogger.Config

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

		<-signals

		signal.Stop(signals)
		cancel()
	}()

	eventLogger := eventlogger.NewEventLogger(
		log,
		&config,
		eventsStreaming,
		repository.NewMongoEventRepository(db),
	)
	eventLogger.Listen()

	<-ctx.Done()
}
