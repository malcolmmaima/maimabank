package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/malcolmmaima/maimabank/db/sqlc"
	"github.com/malcolmmaima/maimabank/token"
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

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// if fromAccount is equal to toAccount, return error
	if fromAccount.ID == req.ToAccountID {
		err := errors.New("from account cannot be equal to to account")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
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

// multi currency transfer e.g. KES to USD by getting exchange rate and converting to USD before transfer is done
func (server *Server) createMultiCurrencyTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// if fromAccount is equal to toAccount, return error
	if fromAccount.ID == req.ToAccountID {
		err := errors.New("from account cannot be equal to to account")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	toAccount, err := server.store.GetAccount(ctx, req.ToAccountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if fromAccount currency is not equal to toAccount currency then convert to target currency before transfer

	if fromAccount.Currency != toAccount.Currency {
		// base currency and target currency should not be the same
		if req.Currency == toAccount.Currency {
			err := errors.New("base currency cannot be the same as target currency")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// get exchange rate
		arg := db.GetExchangeRateParams{
			BaseCurrency:   req.Currency,
			TargetCurrency: toAccount.Currency,
		}

		exchangeRate, err := server.store.GetExchangeRate(ctx, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				err := errors.New("exchange rate not found")
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}

			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// convert amount to target currency... exchangeRate.ExchangeRate is string, convert to int64
		exchangeRateAmount, err := strconv.ParseFloat(exchangeRate.ExchangeRate, 64)
		amountToConvert := float64(req.Amount)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		amount :=  amountToConvert * exchangeRateAmount

		// if amount is less than 1, return error
		if amount < 1 {
			increaseAmountTo := strconv.FormatInt(int64(1 / exchangeRateAmount) + 1, 10)
			err := fmt.Errorf("amount to transfer is less than 1 %s, please increase amount to %s %s or more", toAccount.Currency, req.Currency, increaseAmountTo)
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		arg2 := db.TransferTxParams{
			FromAccountID: req.FromAccountID,
			ToAccountID:   req.ToAccountID,
			Amount:        int64(amount),
		}

		result, err := server.store.TransferTx(ctx, arg2)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, result)
		return
	}

	// if fromAccount currency is equal to toAccount currency then transfer as usual
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


func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch, expecting %s instead of %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
