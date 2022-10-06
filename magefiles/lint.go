package main

func LintGo() error {
	return GoLangCiLint(
		"run",
		"-v",
		"--issues-exit-code=1",
	)
}
