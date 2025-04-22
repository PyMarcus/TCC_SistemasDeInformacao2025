package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/config"
	adapters "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/http"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/usecase"
)

func TestMain(m *testing.M) {
	_, err := config.LoadConfig("../.env")

	if err != nil{
		println("[-] Error to load config " + err.Error())
		os.Exit(1)
	}
	clientService := adapters.NewApiRequestService() 
	clientUsecase := usecase.NewAPIRequestUsecase(clientService)
	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	response, err := clientUsecase.Fetch("https://example.com", header, "") // or post

	if err != nil || response == nil {
		panic("error: " + err.Error())
	}
	if response.StatusCode == http.StatusBadRequest{
		panic("BadRequest error")
	}

	

	os.Exit(m.Run())
}
