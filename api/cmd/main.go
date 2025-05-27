package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/ronnyp07/SportStream/api/internal/app"
	"github.com/ronnyp07/SportStream/api/internal/metrics"
	"github.com/ronnyp07/SportStream/api/internal/pkg/infaestructure/log"

	"emperror.dev/errors"
	config "github.com/ronnyp07/SportStream/api/internal/pkg/config"
)

var Version = "unknown"

// @title SportStream API
// @version 1.0
// @description This is a sample server for SportStream.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email ronny.morel07@gmail.com

// @host localhost:8080
// @BasePath /api/v1
func main() {

	metrics.NewAppInfoMetricsHandler(metrics.Host, Version, runtime.Version())
	ctx := context.Background()

	err := config.Load()
	if err != nil {
		panic(errors.Wrap(err, "loading config"))
	}

	err = log.SetupLogger("Apiservice")
	if err != nil {
		panic(errors.Wrap(err, "setting log"))
	}

	a := app.App{}
	if err := a.Start(ctx); err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("error starting app %v", err))
		os.Exit(1)
	}
}
