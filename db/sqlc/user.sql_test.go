package db

import (
	"context"
	"database/sql"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	// create a random account
	user := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.HashedPassword, user2.HashedPassword)
	require.Equal(t, user.FullName, user2.FullName)
	require.Equal(t, user.Email, user2.Email)

	require.WithinDuration(t, user.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user.CreatedAt, user2.CreatedAt, time.Second)

}

// Testing UpdateUser by updating the user's full name
func TestUpdateUsersFullName(t *testing.T) {
	// create a random account
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	arg := UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{String: newFullName, Valid: true},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqualf(t, oldUser.FullName, updatedUser.FullName, "Full name should be different")

	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, arg.FullName.String, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)

	require.WithinDuration(t, oldUser.PasswordChangedAt, updatedUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, oldUser.CreatedAt, updatedUser.CreatedAt, time.Second)
}

// Testing UpdateUser by updating the user's email
func TestUpdateUsersEmail(t *testing.T) {
	// create a random account
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()
	arg := UpdateUserParams{
		Username: oldUser.Username,
		Email:    sql.NullString{String: newEmail, Valid: true},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqualf(t, oldUser.Email, updatedUser.Email, "Email should be different")

	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, arg.Email.String, updatedUser.Email)

	require.WithinDuration(t, oldUser.PasswordChangedAt, updatedUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, oldUser.CreatedAt, updatedUser.CreatedAt, time.Second)
}

// Testing UpdateUser by updating the user's password
func TestUpdateUsersPassword(t *testing.T) {
	// create a random account
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(6)
	arg := UpdateUserParams{
		Username:       oldUser.Username,
		HashedPassword: sql.NullString{String: newPassword, Valid: true},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqualf(t, oldUser.HashedPassword, updatedUser.HashedPassword, "Password should be different")

	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, arg.HashedPassword.String, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)

	require.WithinDuration(t, oldUser.PasswordChangedAt, updatedUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, oldUser.CreatedAt, updatedUser.CreatedAt, time.Second)
}

// Testing UpdateUser by updating the user's full name, email, and password
func TestUpdateUsersFullNameEmailPassword(t *testing.T) {
	// create a random account
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)
	arg := UpdateUserParams{
		Username:       oldUser.Username,
		FullName:       sql.NullString{String: newFullName, Valid: true},
		Email:          sql.NullString{String: newEmail, Valid: true},
		HashedPassword: sql.NullString{String: newPassword, Valid: true},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqualf(t, oldUser.FullName, updatedUser.FullName, "Full name should be different")
	require.NotEqualf(t, oldUser.Email, updatedUser.Email, "Email should be different")
	require.NotEqualf(t, oldUser.HashedPassword, updatedUser.HashedPassword, "Password should be different")

	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, arg.HashedPassword.String, updatedUser.HashedPassword)
	require.Equal(t, arg.FullName.String, updatedUser.FullName)
	require.Equal(t, arg.Email.String, updatedUser.Email)

	require.WithinDuration(t, oldUser.PasswordChangedAt, updatedUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, oldUser.CreatedAt, updatedUser.CreatedAt, time.Second)
}
