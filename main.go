package main

import (
	"github.com/baez90/inetmock/internal/cmd"
	"go.uber.org/zap"
	"os"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := cmd.ExecuteRootCommand(); err != nil {
		logger.Error("Failed to run inetmock",
			zap.Error(err),
		)
		os.Exit(1)
	}
}
