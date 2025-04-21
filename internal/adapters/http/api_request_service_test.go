package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApiRequestService(t *testing.T){
	clientService := NewApiRequestService()

	assert.NotNil(t, clientService)
}
