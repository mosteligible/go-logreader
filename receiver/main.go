package main

import (
	"github.com/mosteligible/go-logreader/receiver/app"
)

func main() {
	app := app.NewApp()
	app.Run()
}
