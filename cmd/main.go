package main

import (
	"log"

	"github.com/Blue-Onion/MahilAi/handler/config"
	"github.com/Blue-Onion/MahilAi/handler/config/Camera"
)


func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print("Config loaded. Starting camera work...")

	camera.StartCameraWork(cfg)
}