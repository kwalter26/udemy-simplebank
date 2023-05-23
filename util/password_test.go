package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

// TestHashPassword tests the HashPassword function
func TestHashPassword(t *testing.T) {
	password := RandomString(6)
	hashPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	err = CheckPassword(password, hashPassword)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashPassword)
	require.EqualError(t, bcrypt.ErrMismatchedHashAndPassword, err.Error())

	hashPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword2)
	require.NotEqual(t, hashPassword, hashPassword2)

	longPassword := RandomString(1000)
	_, err = HashPassword(longPassword)
	require.EqualError(t, err, bcrypt.ErrPasswordTooLong.Error())

}
