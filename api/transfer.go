package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"net/http"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func (s *Server) createTransfer(context *gin.Context) {
	var req createTransferRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, errorResponse(err))
		return
	}

	if !s.validAccountCurrency(context, req.FromAccountID, req.Currency) {
		return
	}

	if !s.validAccountCurrency(context, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := s.store.TransferTx(context, arg)
	if err != nil {
		context.JSON(500, errorResponse(err))
		return
	}

	context.JSON(200, transfer)
}

// check valid account currency
func (s *Server) validAccountCurrency(context *gin.Context, accountID int64, currency string) bool {
	account, err := s.store.GetAccount(context, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(404, errorResponse(err))
			return false
		}
		context.JSON(500, errorResponse(err))
		return false
	}
	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
