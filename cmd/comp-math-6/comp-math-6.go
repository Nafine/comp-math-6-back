package main

import (
	"comp-math-6/internal/config"
	"comp-math-6/internal/web"
)

func main() {
	cfg, err := config.Get()

	if err != nil {
		panic(err)
	}

	server := web.New(cfg)

	_ = server.Start()
}
