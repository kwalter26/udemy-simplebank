package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Renew Access Token API using Refresh Token

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// renewAccessToken generates a new access token using a refresh token.
// It returns an error if the refresh token is invalid or expired.
func (s *Server) renewAccessToken(context *gin.Context) {
	var req renewAccessTokenRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify refresh token
	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Get the session from the database
	session, err := s.store.GetSession(context, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check if session is blocked
	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check if session username matches
	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check if session token matches request token
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatch session token")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check if session token is expired
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("refresh token expired")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Create a new access token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(refreshPayload.Username, s.config.AccessTokenDuration)
	if err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Return the access token
	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpireAt,
	}

	context.JSON(http.StatusOK, rsp)
}
