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

	alert := app.NewAlert(*configFile)
	if err := alert.Run(); err != nil {
		os.Exit(1)
	}
}
