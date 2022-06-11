package main

import (
	"flag"
	"github.com/zhjx922/alert/app"
	"os"
)

func main() {
	configFile := flag.String("c", "", "Config file path")
	flag.Parse()

	if *configFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	app := app.NewAlert(*configFile)
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
