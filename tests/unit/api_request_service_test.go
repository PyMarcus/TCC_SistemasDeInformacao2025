package unit

import (
	"testing"

	adapters "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/http"
	"github.com/stretchr/testify/assert"
)

func TestNewApiRequestService(t *testing.T){
	clientService := adapters.NewApiRequestService()

	assert.NotNil(t, clientService)
}

func TestSetters(t *testing.T){
	clientService := adapters.NewApiRequestService()

	result := clientService.SetMethod("GET")

	assert.NotEmpty(t, result)

	result = clientService.SetURL("https://www.google.com")

	assert.NotEmpty(t, result) 

	result = clientService.SetBody("")

	assert.NotEmpty(t, result) 

	header := make(map[string]string)
	header["hello"] = "world"
	result = clientService.SetHeaders(header)

	assert.NotEmpty(t, result) 
	
}