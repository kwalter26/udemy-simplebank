package api

import (
	mockdb "github.com/kwalter26/udemy-simplebank/db/mock"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewServerWithNewRelic(t *testing.T) {
	store := mockdb.NewMockStore(nil)
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
		NewRelicAppEnabled:  true,
		NewRelicLicenseKey:  "1234567890123456789012345678901234567890",
		NewRelicAppName:     "test",
	}

	_, err := NewServer(config, store)
	require.NoError(t, err)
}

func TestNewServerWithNewRelicEmptyLicense(t *testing.T) {
	store := mockdb.NewMockStore(nil)
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
		NewRelicAppEnabled:  true,
	}

	_, err := NewServer(config, store)
	require.Error(t, err)
}

func TestNewServerWithNewRelicEmptyName(t *testing.T) {
	store := mockdb.NewMockStore(nil)
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
		NewRelicAppEnabled:  true,
		NewRelicLicenseKey:  "1234567890123456789012345678901234567890",
	}

	_, err := NewServer(config, store)
	require.Error(t, err)
}
