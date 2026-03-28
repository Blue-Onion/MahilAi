package main

import (
	"log"

	camera "github.com/Blue-Onion/MahilAi/handler/Camera"
	"github.com/Blue-Onion/MahilAi/handler/config"
)


func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print("Config loaded. Starting camera work...")

	camera.StartCameraWork(cfg)
}