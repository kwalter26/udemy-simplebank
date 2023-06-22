package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	mockdb "github.com/kwalter26/udemy-simplebank/db/mock"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/pb"
	"github.com/kwalter26/udemy-simplebank/token"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func TestUpdateUserAPI(t *testing.T) {
	user, _ := createRandomUser(t)
	otherUser, _ := createRandomUser(t)

	newName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)

	invalidEmail := "invalid_email"
	invalidFullName := "invalid-full_name"
	invalidPassword := "short"

	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: sql.NullString{
						String: newName,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
				}
				updatedUser := user
				updatedUser.FullName = newName
				updatedUser.Email = newEmail
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				accessToken, _, err := tokenMaker.CreateToken(user.Username, time.Minute)
				require.NoError(t, err)
				require.NotEmpty(t, accessToken)
				bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
				md := metadata.MD{
					authorizationHeader: []string{
						bearerToken,
					},
				}
				return metadata.NewIncomingContext(context.Background(), md)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updatedUser := res.User
				require.Equal(t, user.Username, updatedUser.Username)
				require.Equal(t, newName, updatedUser.FullName)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
		{
			name: "UnauthenticatedUser",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.ErrorContains(t, err, "missing context metadata")
			},
		},
		{
			name: "EmptyUsername",
			req: &pb.UpdateUserRequest{
				Username: "",
			},
			buildStubs: func(store *mockdb.MockStore) {},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return getAuthCtx(t, tokenMaker, user, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.ErrorContains(t, err, "InvalidArgument")
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				Email:    &invalidEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return getAuthCtx(t, tokenMaker, user, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.ErrorContains(t, err, "InvalidArgument")
			},
		},
		{
			name: "InvalidFullName",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &invalidFullName,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return getAuthCtx(t, tokenMaker, user, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.ErrorContains(t, err, "InvalidArgument")
			},
		},
		{
			name: "InvalidPassword",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				Password: &invalidPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return getAuthCtx(t, tokenMaker, user, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.ErrorContains(t, err, "InvalidArgument")
			},
		},
		{
			name: "WrongUser",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return getAuthCtx(t, tokenMaker, otherUser, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.ErrorContains(t, err, "PermissionDenied")
			},
		},
		{
			name: "UserNotFound",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(db.User{}, sql.ErrNoRows)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return getAuthCtx(t, tokenMaker, user, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.ErrorContains(t, err, "NotFound")
			},
		},
		{
			name: "InternalError",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(db.User{}, sql.ErrConnDone)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return getAuthCtx(t, tokenMaker, user, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.ErrorContains(t, err, "Internal")
			},
		},
		{
			name: "ExpiredAccessToken",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return getAuthCtx(t, tokenMaker, user, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.ErrorContains(t, err, "Unauthenticated")
			},
		},
		{
			name: "UpdatePassword",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				Password: &newPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(user, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return getAuthCtx(t, tokenMaker, user, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, res.User.Username, user.Username)
				require.Equal(t, res.User.FullName, user.FullName)
				require.Equal(t, res.User.Email, user.Email)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)

			server := newTestServer(t, store, nil)

			ctx := tc.buildContext(t, server.tokenMaker)
			res, err := server.UpdateUser(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func getAuthCtx(t *testing.T, tokenMaker token.Maker, user db.User, expire time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(user.Username, expire)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}
	return metadata.NewIncomingContext(context.Background(), md)
}
