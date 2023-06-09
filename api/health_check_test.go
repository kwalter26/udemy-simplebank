package api

import (
	"errors"
	"github.com/golang/mock/gomock"
	mockdb "github.com/kwalter26/udemy-simplebank/db/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHealthAndReady(t *testing.T) {

	testCases := []struct {
		name          string
		url           string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			url:  "/healthz",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().Ping(gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, `{"status":"Ok"}`, recorder.Body.String())
			},
		},
		{
			name: "DB Error",
			url:  "/healthz",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().Ping(gomock.Any()).Times(1).Return(errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				require.Equal(t, `{"error":"db error"}`, recorder.Body.String())
			},
		},
		{
			name: "OK",
			url:  "/readyz",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().Ping(gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, `{"status":"Ok"}`, recorder.Body.String())
			},
		},
		{
			name: "DB Error",
			url:  "/readyz",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().Ping(gomock.Any()).Times(1).Return(errors.New("db error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				require.Equal(t, `{"error":"db error"}`, recorder.Body.String())
			},
		},
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

			request, err := http.NewRequest(http.MethodGet, tc.url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
