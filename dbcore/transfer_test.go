package dbcore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createMockTransaction(t *testing.T) (Account, Account, Transfer) {
	accountFrom := createMockAccount(t)
	accountTo := createMockAccount(t)
	args := CreateTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        100,
	}
	transfer, err := testQUeries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.FromAccountID, accountFrom.ID)
	require.Equal(t, transfer.ToAccountID, accountTo.ID)
	return accountFrom, accountTo, transfer
}

func TestCreateTransaction(t *testing.T) {
	createMockTransaction(t)
}

func TestGetTransaction(t *testing.T) {
	_, _, transaction := createMockTransaction(t)
	trans, err := testQUeries.GetTransfer(context.Background(), transaction.ID)
	require.NoError(t, err)
	require.NotEmpty(t, trans)
	require.Equal(t, trans.FromAccountID, transaction.FromAccountID)
	require.Equal(t, trans.ToAccountID, transaction.ToAccountID)
}

func TestUpdateTransfer(t *testing.T) {
	_, _, transfer := createMockTransaction(t)
	args := UpdateTransferParams{
		ID:            transfer.ID,
		FromAccountID: transfer.FromAccountID,
		ToAccountID:   transfer.ToAccountID,
		Amount:        11300,
	}
	trans, err := testQUeries.UpdateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, trans)
	require.Equal(t, trans.FromAccountID, args.FromAccountID)
	require.Equal(t, trans.ToAccountID, args.ToAccountID)
	require.Equal(t, trans.Amount, args.Amount)
}

func TestGetTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createMockTransaction(t)
	}
	args := GetTransfersParams{
		Limit:  5,
		Offset: 5,
	}
	transfers, err := testQUeries.GetTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, trans := range transfers {
		require.NotEmpty(t, trans)
	}
}

func TestDeleteTransfer(t *testing.T) {
	_, _, transfer := createMockTransaction(t)
	err := testQUeries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	trans, err := testQUeries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.Empty(t, trans)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
