package main

import (
	"os"

	"github.com/ardizariat/generate-pdf/internal/cli"
	"github.com/ardizariat/generate-pdf/internal/logging"
	"go.uber.org/zap"
)

func main() {
	logger := logging.New(false)

	if err := cli.New(logger).Execute(); err != nil {
		logger.Error("generate failed", zap.Error(err))
		_ = logger.Sync() // stderr sync often returns EINVAL; not fatal
		os.Exit(1)
	}

	_ = logger.Sync()
}
