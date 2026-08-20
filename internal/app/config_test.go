package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("DB_SOURCE", "test-source")
	t.Setenv("SERVER_ADDRESS", ":9999")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.Equal(t, "test-source", cfg.DBSource)
	require.Equal(t, ":9999", cfg.ServerAddress)
}

func TestLoadConfigRejectsEmptyDBSource(t *testing.T) {
	t.Setenv("DB_SOURCE", "")

	_, err := LoadConfig()
	require.Error(t, err)
}
