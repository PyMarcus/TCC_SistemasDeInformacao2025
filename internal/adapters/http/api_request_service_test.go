package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApiRequestService(t *testing.T){
	clientService := NewApiRequestService()

	assert.NotNil(t, clientService)
}

func TestSetters(t *testing.T){
	clientService := NewApiRequestService()

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