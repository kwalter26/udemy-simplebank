package db

import (
	"context"
	"github.com/kwalter26/udemy-simplebank/token"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Test session.sql.go
func TestCreateSession(t *testing.T) {

	// Create random session
	createRandomSession(t)
}

// Test get session
func TestGetSession(t *testing.T) {

	// Create random session
	session1 := createRandomSession(t)

	// Get session
	session2, err := testQueries.GetSession(context.Background(), session1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	// Compare session1 and session2
	require.Equal(t, session1.ID, session2.ID)
	require.Equal(t, session1.Username, session2.Username)
	require.Equal(t, session1.RefreshToken, session2.RefreshToken)
	require.Equal(t, session1.UserAgent, session2.UserAgent)
	require.Equal(t, session1.ClientIp, session2.ClientIp)
	require.Equal(t, session1.IsBlocked, session2.IsBlocked)
	require.WithinDuration(t, session1.ExpiresAt, session2.ExpiresAt, 4*time.Minute+time.Second)
	require.WithinDuration(t, session1.CreatedAt, session2.CreatedAt, time.Second)
}

// function for createRandomSession
func createRandomSession(t *testing.T) Session {
	randomUser := createRandomUser(t)

	// Create refreshToken
	maker, err := token.NewPasetoMaker(util.RandomString(32))

	refreshToken, p, err := maker.CreateToken(randomUser.Username, 4*time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, refreshToken)
	require.NotEmpty(t, p)

	// Create new db createSession object to test with
	createSession := CreateSessionParams{
		ID:           p.ID,
		Username:     randomUser.Username,
		RefreshToken: refreshToken,
		UserAgent:    "userAgent",
		ClientIp:     "clientIp",
		IsBlocked:    false,
		ExpiresAt:    p.ExpireAt,
	}

	// Create session in db
	session, err := testQueries.CreateSession(context.Background(), createSession)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, createSession.ID, session.ID)
	require.Equal(t, createSession.Username, session.Username)
	require.Equal(t, createSession.RefreshToken, session.RefreshToken)
	require.Equal(t, createSession.UserAgent, session.UserAgent)
	require.Equal(t, createSession.ClientIp, session.ClientIp)
	require.Equal(t, createSession.IsBlocked, session.IsBlocked)
	require.WithinDuration(t, createSession.ExpiresAt, session.ExpiresAt, 4*time.Minute+time.Second)
	require.WithinDuration(t, time.Now(), session.CreatedAt, time.Second)

	return session
}
