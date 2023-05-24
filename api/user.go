package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/lib/pq"
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
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (s *Server) loginUser(context *gin.Context) {
	var req LoginUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(context, req.Username)
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

	accessToken, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		context.JSON(500, errorResponse(err))
		return
	}

	rsp := LoginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	context.JSON(200, rsp)
}
