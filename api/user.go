package api

import (
	ctx "context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/lib/pq"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
	"net/http"
	"time"
)

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=40"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=40"`
	Email    string `json:"email" binding:"required,email"`
}

// CreateUserResponse represents a response from a create user request
type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// CreateUser creates a new user account
func (s *Server) CreateUser(context *gin.Context) {
	var req CreateUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(context, arg)
	if err != nil {
		if pgErr, err := err.(*pq.Error); err {
			switch pgErr.Code.Name() {
			case "unique_violation":
				context.JSON(400, errorResponse(pgErr))
				return
			}
		}
		context.JSON(500, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)

	context.JSON(200, rsp)
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=40"`
	Password string `json:"password" binding:"required,min=6,max=40"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (s *Server) loginUser(context *gin.Context) {
	var req LoginUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	txn := s.app.StartTransaction("postgresQuery")
	log.Println("txn start", txn)
	c := newrelic.NewContext(ctx.Background(), txn)
	user, err := s.store.GetUser(c, req.Username)
	txn.End()
	log.Println("txn end", txn)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		context.JSON(500, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := s.store.CreateSession(context, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    context.Request.UserAgent(),
		ClientIp:     context.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpireAt,
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpireAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpireAt,
		User:                  newUserResponse(user),
	}
	context.JSON(200, rsp)
}
