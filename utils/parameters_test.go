package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadParameters(t *testing.T) {
	t.Run("Load parameters yaml file", func(t *testing.T) {
		os.Setenv("config", "../config/test/parameters.yaml")
		parameters := LoadParameters()

		assert.Equal(t, 8080, parameters.Port)
		assert.Equal(t, "./storage", parameters.StorageRoot)
	})
}
