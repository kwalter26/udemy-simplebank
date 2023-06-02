package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/kwalter26/udemy-simplebank/db/mock"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/token"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRenewToken(t *testing.T) {

	user, _ := createRandomUser(t)

	testCases := []struct {
		name          string
		body          func(s *Server) (string, *token.Payload)
		buildStubs    func(store *mockdb.MockStore, token string, payload *token.Payload)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: func(s *Server) (string, *token.Payload) {
				validToken, payload, _ := s.tokenMaker.CreateToken(user.Username, time.Minute)
				return validToken, payload
			},
			buildStubs: func(store *mockdb.MockStore, token string, payload *token.Payload) {
				arg := db.Session{
					ID:           payload.ID,
					Username:     payload.Username,
					ExpiresAt:    payload.ExpireAt,
					RefreshToken: token,
					IsBlocked:    false,
				}
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(payload.ID)).Times(1).Return(arg, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.NotEmpty(t, recorder.Body)
			},
		},
		{
			name: "ExpiredRefreshToken",
			body: func(s *Server) (string, *token.Payload) {
				validToken, payload, _ := s.tokenMaker.CreateToken(user.Username, -time.Minute)
				return validToken, payload
			},
			buildStubs: func(store *mockdb.MockStore, token string, payload *token.Payload) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: func(s *Server) (string, *token.Payload) {
				validToken, payload, _ := s.tokenMaker.CreateToken(user.Username, time.Minute)
				return validToken, payload
			},
			buildStubs: func(store *mockdb.MockStore, token string, payload *token.Payload) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(payload.ID)).Times(1).Return(db.Session{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: func(s *Server) (string, *token.Payload) {
				validToken, payload, _ := s.tokenMaker.CreateToken(user.Username, time.Minute)
				return validToken, payload
			},
			buildStubs: func(store *mockdb.MockStore, token string, payload *token.Payload) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(payload.ID)).Times(1).Return(db.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidRefreshToken",
			body: func(s *Server) (string, *token.Payload) {
				return "asdfasdf", nil
			},
			buildStubs: func(store *mockdb.MockStore, token string, payload *token.Payload) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "IsBlocked",
			body: func(s *Server) (string, *token.Payload) {
				validToken, payload, _ := s.tokenMaker.CreateToken(user.Username, time.Minute)
				return validToken, payload
			},
			buildStubs: func(store *mockdb.MockStore, token string, payload *token.Payload) {
				arg := db.Session{
					ID:           payload.ID,
					Username:     payload.Username,
					ExpiresAt:    payload.ExpireAt,
					RefreshToken: token,
					IsBlocked:    true,
				}
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(payload.ID)).Times(1).Return(arg, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Wrong Username",
			body: func(s *Server) (string, *token.Payload) {
				validToken, payload, _ := s.tokenMaker.CreateToken("wrong", time.Minute)
				return validToken, payload
			},
			buildStubs: func(store *mockdb.MockStore, token string, payload *token.Payload) {
				arg := db.Session{
					ID:           payload.ID,
					Username:     "wrong wrong",
					ExpiresAt:    payload.ExpireAt,
					RefreshToken: token,
					IsBlocked:    false,
				}
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(payload.ID)).Times(1).Return(arg, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Wrong token",
			body: func(s *Server) (string, *token.Payload) {
				validToken, payload, _ := s.tokenMaker.CreateToken("wrong", time.Minute)
				return validToken, payload
			},
			buildStubs: func(store *mockdb.MockStore, token string, payload *token.Payload) {
				arg := db.Session{
					ID:           payload.ID,
					Username:     payload.Username,
					ExpiresAt:    payload.ExpireAt,
					RefreshToken: "token",
					IsBlocked:    false,
				}
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(payload.ID)).Times(1).Return(arg, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Expired session",
			body: func(s *Server) (string, *token.Payload) {
				validToken, payload, _ := s.tokenMaker.CreateToken("wrong", time.Minute)
				return validToken, payload
			},
			buildStubs: func(store *mockdb.MockStore, token string, payload *token.Payload) {
				arg := db.Session{
					ID:           payload.ID,
					Username:     payload.Username,
					ExpiresAt:    time.Now().Add(-time.Hour),
					RefreshToken: token,
					IsBlocked:    false,
				}
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(payload.ID)).Times(1).Return(arg, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			refreshToken, payload := tc.body(server)
			tc.buildStubs(store, refreshToken, payload)

			url := "/token/renew_access"

			// Create a new refresh token request body
			// and encode it to JSON
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(gin.H{
				"refresh_token": refreshToken,
			})
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, &buf)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}
