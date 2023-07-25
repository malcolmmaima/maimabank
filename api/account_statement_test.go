// account_statement_test.go

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/malcolmmaima/maimabank/db/mock"
	db "github.com/malcolmmaima/maimabank/db/sqlc"
	"github.com/malcolmmaima/maimabank/token"
	"github.com/stretchr/testify/require"
)

func TestListTransfersAPI(t *testing.T) {
	const pageSize = 10

	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	account1 := createRandomAccount(user1.Username)
	account2 := createRandomAccount(user2.Username)

	transfer1 := db.Transfer{
		ID:            1,
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        100,
		CreatedAt:     time.Now(),
	}

	transfer2 := db.Transfer{
		ID:            2,
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        200,
		CreatedAt:     time.Now(),
	}

	testCases := []struct {
		name          string
		authUsername  string
		queryParams   string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			authUsername: user1.Username,
			queryParams:  "account_id=" + strconv.FormatInt(account1.ID, 10) + "&page_id=1&page_size=" + strconv.FormatInt(pageSize, 10),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				listTransfersParams := db.ListTransfersParams{
					FromAccountID: account1.ID,
					ToAccountID:   account1.ID,
					Limit:         pageSize,
					Offset:        0,
				}

				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().ListTransfers(gomock.Any(), gomock.Eq(listTransfersParams)).Times(1).Return([]db.Transfer{transfer1, transfer2}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var responseTransfers []db.Transfer
				err := json.Unmarshal(recorder.Body.Bytes(), &responseTransfers)
				require.NoError(t, err)

				require.Len(t, responseTransfers, 2)
				require.Equal(t, -transfer1.Amount, responseTransfers[0].Amount) 
				require.Equal(t, -transfer2.Amount, responseTransfers[1].Amount)
			},
		},
		// Add more test cases as needed to cover different scenarios
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/accounts/statement?" + tc.queryParams
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
