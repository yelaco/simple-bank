package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/yelaco/simple-bank/db/sqlc"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.validateAccount(ctx, req.FromAccountID, req.Currency); err != nil {
		ctx.Errors = append(ctx.Errors, &gin.Error{
			Err:  fmt.Errorf("api.Server.createTransfer: invalid 'from' account: %w", err),
			Type: gin.ErrorTypePublic,
		})
		return
	}

	if err := server.validateAccount(ctx, req.ToAccountID, req.Currency); err != nil {
		ctx.Errors = append(ctx.Errors, &gin.Error{
			Err:  fmt.Errorf("api.Server.createTransfer: invalid 'to' account: %w", err),
			Type: gin.ErrorTypePublic,
		})
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) error {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return fmt.Errorf("api.Server.validateAccount: account not found: %w", err)
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return err
	}

	if account.Currency != currency {
		err := fmt.Errorf("api.Server.validateAccount: account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return err
	}

	return nil
}
