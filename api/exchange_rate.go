package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/malcolmmaima/maimabank/db/sqlc"
	"github.com/malcolmmaima/maimabank/token"
	"github.com/malcolmmaima/maimabank/util"
)

// get exchange rate by passing id of currency
type exchangeRateRequest struct {
	CurrencyID int64 `form:"currency_id" binding:"required"`
}

func (server *Server) getExchangeRate(ctx *gin.Context) {
	var req exchangeRateRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get exchange rate
	exchangeRate, err := server.store.GetExchangeRate(ctx, req.CurrencyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err == sql.ErrNoRows {
		err := errors.New("currency not found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	// Return exchange rate
	ctx.JSON(http.StatusOK, exchangeRate)
}


// create exchange rate
type createExchangeRateRequest struct {
	BaseCurrency   string `json:"base_currency" binding:"required,currency"`
	TargetCurrency string `json:"target_currency" binding:"required,currency"`
	ExchangeRate   string `json:"exchange_rate" binding:"required,gt=0"`
}

func (server *Server) createExchangeRate(ctx *gin.Context) {
	var req createExchangeRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateExchangeRateParams{
		BaseCurrency:   req.BaseCurrency,
		TargetCurrency: req.TargetCurrency,
		ExchangeRate:   req.ExchangeRate,
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != "admin" {
		err := errors.New("only admin can create exchange rates")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// make sure base currency and target currency are not the same and are valid/supported
	if req.BaseCurrency == req.TargetCurrency {
		err := errors.New("base currency cannot be the same as target currency")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// make sure base currency and target currency are valid/supported
	validCurrency := util.IsSupportedCurrency(req.BaseCurrency)
	if !validCurrency {
		err := errors.New("base currency not supported")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	validCurrency = util.IsSupportedCurrency(req.TargetCurrency)
	if !validCurrency {
		err := errors.New("target currency not supported")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Create exchange rate
	exchangeRate, err := server.store.CreateExchangeRate(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				err := errors.New("exchange rate already exists")
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Return exchange rate
	ctx.JSON(http.StatusOK, exchangeRate)
}