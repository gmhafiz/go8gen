package main

import (
	"log"

	"abc/configs"
	"abc/internal/server"
)

const Version = "v0.1.0"

func main() {
	cfg := configs.New()

	s := server.NewApp(cfg)

	if err := s.Run(cfg, Version); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
