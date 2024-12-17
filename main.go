package main

import (
	"chat/cmd/app"
	"chat/cmd/migration"
	"flag"
)

func main() {
	//flag
	var cmd string
	flag.StringVar(&cmd, "cmd", "app", "command to run")

	flag.Parse()
	switch cmd {
	case "app":
		app.Start()
	case "migration":
		migration.Start()
	default:
		app.Start()
	}
}
