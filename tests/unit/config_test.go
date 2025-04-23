package unit

import (
	"testing"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T){
	cfg, err := config.LoadConfig("../../.env")

	assert.NoError(t, err)
	assert.NotEmpty(t, cfg.DatabaseUrl)
}