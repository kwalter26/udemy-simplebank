package db

import (
	"context"
	"database/sql"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// private func for creating random transfer
func createRandomTransfer(t *testing.T) (Transfer, Account, Account) {
	// create random account
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	args := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}
	transfers, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Equal(t, args.FromAccountID, transfers.FromAccountID)
	require.Equal(t, args.ToAccountID, transfers.ToAccountID)
	require.Equal(t, args.Amount, transfers.Amount)
	require.NotZero(t, transfers.ID)
	require.NotZero(t, transfers.CreatedAt)
	return transfers, account1, account2
}

// TestCreateTransfer test create transfer
func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

// TestGetTransfer test get transfer
func TestGetTransfer(t *testing.T) {
	transfers1, _, _ := createRandomTransfer(t)
	transfers2, err := testQueries.GetTransfer(context.Background(), transfers1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfers2)
	require.Equal(t, transfers1.ID, transfers2.ID)
	require.Equal(t, transfers1.FromAccountID, transfers2.FromAccountID)
	require.Equal(t, transfers1.ToAccountID, transfers2.ToAccountID)
	require.Equal(t, transfers1.Amount, transfers2.Amount)
	require.WithinDuration(t, transfers1.CreatedAt, transfers2.CreatedAt, time.Second)
}

// TestDeleteTransfer test delete transfer
func TestDeleteTransfer(t *testing.T) {
	transfers1, _, _ := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfers1.ID)
	require.NoError(t, err)
	transfers2, err := testQueries.GetTransfer(context.Background(), transfers1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfers2)
}

// TestListTransfers test list transfers
func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}
	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfers := range transfers {
		require.NotEmpty(t, transfers)
	}
}

// TestUpdateTransfer test update transfer
func TestUpdateTransfer(t *testing.T) {
	transfers1, _, _ := createRandomTransfer(t)
	arg := UpdateTransferParams{
		ID:     transfers1.ID,
		Amount: util.RandomBalance(),
	}
	err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	transfers2, err := testQueries.GetTransfer(context.Background(), transfers1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfers2)
	require.Equal(t, transfers1.ID, transfers2.ID)
	require.Equal(t, arg.Amount, transfers2.Amount)
	require.NotEqual(t, transfers1.Amount, transfers2.Amount)
	require.WithinDuration(t, transfers1.CreatedAt, transfers2.CreatedAt, time.Second)
}
