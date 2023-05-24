package db

import (
	"context"
	"database/sql"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	// create a random account
	account := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account.ID, account2.ID)
	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, account.Balance, account2.Balance)
	require.Equal(t, account.Currency, account2.Currency)
	require.WithinDuration(t, account.CreatedAt, account2.CreatedAt, 0)
}

// Test updating an account
func TestUpdateAccount(t *testing.T) {
	// create a random account
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomBalance(),
	}
	// update the account
	_, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)

	// get the account
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// check if the account is updated
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.NotEqualf(t, account1.Balance, account2.Balance, "account balance should be different")
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 0)

}

// Test deleting an account
func TestDeleteAccount(t *testing.T) {
	// create a random account
	account1 := createRandomAccount(t)

	// delete the account
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	// get the account
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.Empty(t, account2)

	// check if the account is deleted
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

// Test listing accounts
func TestListAccounts(t *testing.T) {
	var lastAccount Account
	// create 10 random accounts
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	// list the accounts
	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	// check if the accounts are listed
	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.ID, account.ID)
	}
}
