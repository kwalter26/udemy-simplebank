package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomOwner(t *testing.T) {
	require.Equal(t, 6, len(RandomOwner()))
}

func TestRandomBalance(t *testing.T) {
	require.GreaterOrEqual(t, int64(1000), RandomBalance())
}

func TestRandomEmail(t *testing.T) {
	require.Contains(t, RandomEmail(), "@")
}

func TestRandomInt(t *testing.T) {
	i := RandomInt(0, 10)
	require.LessOrEqual(t, i, int64(10))
	require.GreaterOrEqual(t, i, int64(0))
}

func TestRandomCurrency(t *testing.T) {
	require.Contains(t, []string{"USD", "CAD", "EUR"}, RandomCurrency())
}
