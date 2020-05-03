package cmd

import (
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/spf13/cobra"
)

var (
	logger    logging.Logger
	serverCmd *cobra.Command

	pluginsDirectory string
	configFilePath   string
	logLevel         string
	developmentLogs  bool
)

func init() {
	serverCmd = &cobra.Command{
		Use:   "",
		Short: "INetMock is lightweight internet mock",
	}

	serverCmd.PersistentFlags().StringVar(&pluginsDirectory, "plugins-directory", "", "Directory where plugins should be loaded from")
	serverCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "Path to config file that should be used")
	serverCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "logging level to use")
	serverCmd.PersistentFlags().BoolVar(&developmentLogs, "development-logs", false, "Enable development mode logs")

	serverCmd.AddCommand(
		serveCmd,
		generateCaCmd,
	)
}

func ExecuteServerCommand() error {
	return serverCmd.Execute()
}

func ExecuteClientCommand() error {
	return cliCmd.Execute()
}
