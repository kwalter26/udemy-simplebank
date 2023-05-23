package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/kwalter26/udemy-simplebank/db/mock"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/util"
	pg "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParamsMatcher(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg: arg, password: password}
}

// TestCreateUserAPI tests the CreateUser API
func TestCreateUserAPI(t *testing.T) {
	randomUser, password := createRandomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  randomUser.Username,
				"full_name": randomUser.FullName,
				"password":  password,
				"email":     randomUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: randomUser.Username,
					FullName: randomUser.FullName,
					Email:    randomUser.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParamsMatcher(arg, password)).
					Times(1).
					Return(randomUser, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, randomUser)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  randomUser.Username,
				"full_name": randomUser.FullName,
				"password":  password,
				"email":     randomUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":  randomUser.Username,
				"full_name": randomUser.FullName,
				"password":  password,
				"email":     randomUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pg.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DuplicateEmail",
			body: gin.H{
				"username":  randomUser.Username,
				"full_name": randomUser.FullName,
				"password":  password,
				"email":     randomUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pg.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "invalid-user#name",
				"full_name": randomUser.FullName,
				"password":  password,
				"email":     randomUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":  randomUser.Username,
				"full_name": randomUser.FullName,
				"password":  "123",
				"email":     randomUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooLongPassword",
			body: gin.H{
				"username":  randomUser.Username,
				"full_name": randomUser.FullName,
				"password":  strings.Repeat("a", 1001),
				"email":     randomUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":  randomUser.Username,
				"full_name": randomUser.FullName,
				"password":  password,
				"email":     "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/users"
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, &buf)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	var createdUser db.User
	err := json.NewDecoder(body).Decode(&createdUser)
	require.NoError(t, err)

	require.NotEmpty(t, createdUser)
	require.Equal(t, user.Username, createdUser.Username)
	require.Equal(t, user.FullName, createdUser.FullName)
	require.Equal(t, user.Email, createdUser.Email)
}

func createRandomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}
	return user, password
}
