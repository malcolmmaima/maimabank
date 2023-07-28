package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/malcolmmaima/maimabank/db/sqlc"
	"github.com/malcolmmaima/maimabank/token"
)

type getAccountTransfersRequest struct {
	AccountID    int64     `form:"account_id" binding:"required,min=1"`
	PageID       int32     `form:"page_id" binding:"required,min=1"`
	PageSize     int32     `form:"page_size" binding:"required,min=5,max=10"`
	StartDate    time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate      time.Time `form:"end_date" time_format:"2006-01-02"`
}

func newTransferResponse(transfers []db.Transfer, userAccountID int64) []db.Transfer {
	var transferResponses []db.Transfer

	for _, transfer := range transfers {
		amount := transfer.Amount
		if transfer.FromAccountID == userAccountID {
			// Treat it as a debit, so append a negative sign to the amount
			amount = -amount
		} else {
			// Treat it as a credit, so keep the amount positive
			amount = transfer.Amount
		}

		// Create the transfer response and append it to the result
		transferResponse := db.Transfer{
			ID:            transfer.ID,
			FromAccountID: transfer.FromAccountID,
			ToAccountID:   transfer.ToAccountID,
			Amount:        amount,
			CreatedAt:     transfer.CreatedAt,
		}
		transferResponses = append(transferResponses, transferResponse)
	}

	return transferResponses
}



// user account statement by ListTransfers query where from_account_id = account_id
func (server *Server) listTransfers(ctx *gin.Context) {
	var req getAccountTransfersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Check if user owns account
	account, err := server.store.GetAccount(ctx, req.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if req.AccountID != account.ID {
			err := errors.New("account does not belong to authenticated user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	listTransfersParams := db.ListTransfersParams{
		FromAccountID: req.AccountID,
		ToAccountID: req.AccountID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	listTransfersByDateParams := db.ListTransfersByDateParams{
		FromAccountID: req.AccountID,
		ToAccountID: req.AccountID,
		CreatedAt: req.StartDate,
		CreatedAt_2: req.EndDate,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,

	}

	var transfers []db.Transfer
	if req.StartDate.IsZero() && req.EndDate.IsZero() {
		transfers, err = server.store.ListTransfers(ctx, listTransfersParams)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} 

	// if either is dates are empty, return error
	if req.StartDate.IsZero() && !req.EndDate.IsZero() || !req.StartDate.IsZero() && req.EndDate.IsZero() {
		err := errors.New("both start_date and end_date must be provided")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// if both dates are provided, use ListTransfersByDate query
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		transfers, err = server.store.ListTransfersByDate(ctx, listTransfersByDateParams)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// if start_date is greater than end_date, return error
	if req.StartDate.After(req.EndDate) {
		err := errors.New("start_date cannot be greater than end_date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// if no transfers are found, return empty array
	if len(transfers) == 0 {
		ctx.JSON(http.StatusOK, []db.Transfer{})
		return
	}

	ctx.JSON(http.StatusOK, newTransferResponse(transfers, req.AccountID))
}


