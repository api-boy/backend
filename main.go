package main

import (
	"apiboy/backend/src/config"
	"apiboy/backend/src/service"
)

func main() {
	// read config
	conf := config.New()

	// setup service
	svc, err := service.New(conf)
	if err != nil {
		panic(err)
	}

	// run service
	svc.Run()
}
