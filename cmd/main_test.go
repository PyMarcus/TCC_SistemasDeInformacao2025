package main

import (
	"os"
	"testing"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/config"
)

func TestMain(m *testing.M) {
	_, err := config.LoadConfig("../.env")

	if err != nil{
		println("[-] Error to load config " + err.Error())
		os.Exit(1)
	}

	os.Exit(m.Run())
}
