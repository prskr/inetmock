package main

func main() {
	defer appCancel()
	if err := cliCmd.Execute(); err != nil {
		panic(err)
	}
}
