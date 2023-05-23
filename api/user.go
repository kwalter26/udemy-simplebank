package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/lib/pq"
)

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=40"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=40"`
	Email    string `json:"email" binding:"required,email"`
}

// CreateUserResponse represents a response from a create user request
type CreateUserResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
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

	rsp := CreateUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}

	context.JSON(200, rsp)
}
