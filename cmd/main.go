package main

import (
	"log"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/config"
)

func main() {
	cfg, err := config.LoadConfig(".env")

	if err != nil{
		log.Fatal(err)
	}

	log.Println(cfg.DatabaseUrl)
}