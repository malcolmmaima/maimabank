package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/malcolmmaima/maimabank/db/sqlc"
	"github.com/malcolmmaima/maimabank/token"
)

type getAccountTransfersRequest struct {
	AccountID int64 `form:"account_id" binding:"required,min=1"`
	PageID  int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func newTransferResponse(transfers []db.Transfer, userAccountID int64) []db.Transfer {
	var transferResponses []db.Transfer

	for _, transfer := range transfers {
		amount := transfer.Amount
		if transfer.FromAccountID == userAccountID {
			// Treat it as a debit, so append a negative sign to the amount
			amount = -amount
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

	transfers, err := server.store.ListTransfers(ctx, listTransfersParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newTransferResponse(transfers, req.AccountID))
}


