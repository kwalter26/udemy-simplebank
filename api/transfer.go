package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/token"
	"net/http"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(context *gin.Context) {
	var req createTransferRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, errorResponse(err))
		return
	}

	fromAccount, valid := s.validAccountCurrency(context, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := fmt.Errorf("from account doesn't belong to the authenticated user")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = s.validAccountCurrency(context, req.ToAccountID, req.Currency)
	if !valid {
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
func (s *Server) validAccountCurrency(context *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(context, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(404, errorResponse(err))
			return account, false
		}
		context.JSON(500, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
