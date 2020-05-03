package main

import "github.com/baez90/inetmock/internal/cmd"

func main() {
	if err := cmd.ExecuteClientCommand(); err != nil {
		panic(err)
	}
}
