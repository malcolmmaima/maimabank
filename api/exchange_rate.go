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

// get exchange rate by passing base currency and target currency
type exchangeRateRequest struct {
	BaseCurrency   string `form:"base_currency" binding:"required,currency"`
	TargetCurrency string `form:"target_currency" binding:"required,currency"`
}

func (server *Server) getExchangeRate(ctx *gin.Context) {
	var req exchangeRateRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {

		if req.BaseCurrency == "" || req.TargetCurrency == "" {
			err := errors.New("base currency and target currency cannot be empty")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

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

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// base currency and target currency should not be the same
	if req.BaseCurrency == req.TargetCurrency {
		err := errors.New("base currency cannot be the same as target currency")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Get exchange rate
	arg := db.GetExchangeRateParams{
		BaseCurrency:   req.BaseCurrency,
		TargetCurrency: req.TargetCurrency,
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

// update exchange rate by passing ID and exchange rate
type updateExchangeRateRequest struct {
	ID int64 `json:"id" binding:"required,min=1" uri:"id"`
	ExchangeRate   string `json:"exchange_rate" binding:"required,gt=0"`
}


func (server *Server) updateExchangeRate(ctx *gin.Context) {
	var req updateExchangeRateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateExchangeRateParams{
		ID:             req.ID,
		ExchangeRate:   req.ExchangeRate,
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != "admin" {
		err := errors.New("only admin can update exchange rates")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// make sure exchange rate is not zero or negative
	if req.ExchangeRate <= "0" {
		err := errors.New("exchange rate cannot be zero or negative")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Update exchange rate
	exchangeRate, err := server.store.UpdateExchangeRate(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
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

// list exchange rates
func (server *Server) listExchangeRates(ctx *gin.Context) {
	exchangeRates, err := server.store.ListExchangeRates(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			err := errors.New("no exchange rates found")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, exchangeRates)
}
