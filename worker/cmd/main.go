package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/ronnyp07/SportStream/worker/internal/app"
	"github.com/ronnyp07/SportStream/worker/internal/metrics"
	"github.com/ronnyp07/SportStream/worker/internal/pkg/infaestructure/log"

	"emperror.dev/errors"
	config "github.com/ronnyp07/SportStream/worker/internal/pkg/config"
)

var Version = "unknown"

// @title Worker APi
// @version 1.0
// @contact.name Ronny Morel
// @contact.email ronny.morel01@gmail.com
// @description
func main() {

	metrics.NewAppInfoMetricsHandler(metrics.Host, Version, runtime.Version())
	ctx := context.Background()

	err := config.Load()
	if err != nil {
		panic(errors.Wrap(err, "loading config"))
	}

	err = log.SetupLogger("Workerservice")
	if err != nil {
		panic(errors.Wrap(err, "setting log"))
	}

	a := app.App{}
	if err := a.Start(ctx); err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("error starting app %v", err))
		os.Exit(1)
	}
}
