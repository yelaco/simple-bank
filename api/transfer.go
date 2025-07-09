package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/yelaco/simple-bank/db/sqlc"
	"github.com/yelaco/simple-bank/token"
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

	fromAccount, err := server.validateAccount(ctx, req.FromAccountID, req.Currency)
	if err != nil {
		ctx.Errors = append(ctx.Errors, &gin.Error{
			Err:  fmt.Errorf("api.Server.createTransfer: invalid 'from' account: %w", err),
			Type: gin.ErrorTypePublic,
		})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("api.Server.createTransfer: from account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, err = server.validateAccount(ctx, req.ToAccountID, req.Currency)
	if err != nil {
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

func (server *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, error) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return db.Account{}, fmt.Errorf("api.Server.validateAccount: account not found: %w", err)
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return db.Account{}, err
	}

	if account.Currency != currency {
		err := fmt.Errorf("api.Server.validateAccount: account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return db.Account{}, err
	}

	return account, nil
}
