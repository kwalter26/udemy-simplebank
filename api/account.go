package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/token"
	"github.com/lib/pq"
	"net/http"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(context *gin.Context) {
	var req createAccountRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, errorResponse(err))
		return
	}

	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
	}

	account, err := s.store.CreateAccount(context, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				context.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		context.JSON(500, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(context *gin.Context) {
	var req getAccountRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(400, errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(context, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err = errors.New("account doesn't belong to the authenticated user")
		context.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccounts(context *gin.Context) {
	var req listAccountsRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(400, errorResponse(err))
		return
	}

	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := s.store.ListAccounts(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, accounts)

}

type updateAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateAccountRequestBody struct {
	Balance int64 `json:"balance" binding:"required,min=0"`
}

func (s *Server) updateAccount(context *gin.Context) {
	var reqUri updateAccountRequest
	var reqBody updateAccountRequestBody
	if err := context.ShouldBindUri(&reqUri); err != nil {
		context.JSON(400, errorResponse(err))
		return
	}

	if err := context.ShouldBindJSON(&reqBody); err != nil {
		context.JSON(400, errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      reqUri.ID,
		Balance: reqBody.Balance,
	}

	account, err := s.store.UpdateAccount(context, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, account)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteAccount(context *gin.Context) {
	var req deleteAccountRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(400, errorResponse(err))
		return
	}

	err := s.store.DeleteAccount(context, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, "Account deleted")
}
