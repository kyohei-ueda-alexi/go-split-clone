package main

import (
	"os"
	"split/cli"
)

func main() {
	exitCode := cli.Split(os.Args)
	os.Exit(exitCode)
}
