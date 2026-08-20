package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Setenv("DB_SOURCE", "test-source")
	t.Setenv("SERVER_ADDRESS", ":9999")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, "test-source", cfg.DBSource)
	require.Equal(t, ":9999", cfg.ServerAddress)
}

func TestLoadRejectsEmptyDBSource(t *testing.T) {
	t.Setenv("DB_SOURCE", "")

	_, err := Load()
	require.Error(t, err)
}
