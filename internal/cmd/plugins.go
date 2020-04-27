package cmd

import (
	"github.com/spf13/cobra"
)

var (
	pluginsCmd = &cobra.Command{
		Use:   "plugins",
		Short: "Use the plugins prefix to interact with commands that are provided by plugins",
		Long: `
The plugin prefix can be used to interact with commands that are provided by plugins.
The easiest way to explore what commands are available is to start with 'inetmock plugins' - like you did!
This help page contains a list of available sub-commands starting with the name of the plugin as a prefix.
`,
	}
)
