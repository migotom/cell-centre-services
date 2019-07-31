package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"go.uber.org/zap"

	"github.com/migotom/cell-centre-services/pkg/services/restapi"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var config restapi.Config

	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend("/etc/cell-centre/restapi/config.toml"),
	)
	if err := loader.Load(context.Background(), &config); err != nil {
		log.Fatal("Error during config reading", zap.Error(err))
	}

	restAPI := restapi.NewRESTAPI(
		log,
		&config,
	)
	restAPI.Listen()
}
