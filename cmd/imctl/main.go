package main

import "gitlab.com/inetmock/inetmock/internal/cmd"

func main() {
	if err := cmd.ExecuteClientCommand(); err != nil {
		panic(err)
	}
}
