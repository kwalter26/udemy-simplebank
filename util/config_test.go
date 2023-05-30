package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// Test func LoadConfig with test.env file
func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("../", Local)
	require.NoError(t, err)
	require.NotEmpty(t, config)
}

func TestLoadConfigNotFound(t *testing.T) {
	config, err := LoadConfig("../../", Local)
	require.NoError(t, err)
	require.NotEmpty(t, config)
}

func TestLoadConfigBadConfig(t *testing.T) {
	_, err := LoadConfig("./", "bad")
	require.Error(t, err)
}
