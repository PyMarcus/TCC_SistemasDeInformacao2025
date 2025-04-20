package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T){
	cfg, err := LoadConfig("")

	assert.NoError(t, err)
	assert.NotEmpty(t, cfg.DatabaseUrl)
}