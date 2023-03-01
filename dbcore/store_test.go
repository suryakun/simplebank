package dbcore

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTransaction(t *testing.T) {

	store := NewStore(dbConn)

	errChan := make(chan error)
	resultChan := make(chan TransferTransactionResult)

	accountA := createMockAccount(t)
	accountB := createMockAccount(t)
	fmt.Println("start >>", accountA.Balance, accountB.Balance)
	amount := int64(10)
	n := 10

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTransaction(context.Background(), TransferTransactionParams{
				FromAccountID: accountA.ID,
				ToAccountID:   accountB.ID,
				Amount:        amount,
			})
			errChan <- err
			resultChan <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)

		result := <-resultChan
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, result.Transfer.FromAccountID, accountA.ID)
		require.Equal(t, result.Transfer.ToAccountID, accountB.ID)
		require.Equal(t, result.Transfer.Amount, amount)

		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, result.FromEntry.AccountID, accountA.ID)
		require.Equal(t, result.FromEntry.Amount, -amount)

		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, result.ToEntry.AccountID, accountB.ID)
		require.Equal(t, result.ToEntry.Amount, amount)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountA.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountB.ID, toAccount.ID)

		fmt.Println(">>> tx", fromAccount.Balance, accountA.Balance, toAccount.Balance, accountB.Balance)

		diff1 := accountA.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - accountB.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccountA, err := testQUeries.GetAccount(context.Background(), accountA.ID)
	require.NoError(t, err)

	updatedAccountB, err := testQUeries.GetAccount(context.Background(), accountB.ID)
	require.NoError(t, err)
	fmt.Println("end >>", updatedAccountA.Balance, updatedAccountB.Balance)
}

func TestTransferTransactionDeadlock(t *testing.T) {

	store := NewStore(dbConn)

	errChan := make(chan error)

	accountA := createMockAccount(t)
	accountB := createMockAccount(t)
	fmt.Println("start >>", accountA.Balance, accountB.Balance)
	amount := int64(10)
	n := 10

	for i := 0; i < n; i++ {
		fromAccountID := accountA.ID
		toAccountID := accountB.ID

		if i%2 == 0 {
			fromAccountID = accountB.ID
			toAccountID = accountA.ID
		}
		go func() {
			_, err := store.TransferTransaction(context.Background(), TransferTransactionParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errChan <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)
	}

	updatedAccountA, err := testQUeries.GetAccount(context.Background(), accountA.ID)
	require.NoError(t, err)

	updatedAccountB, err := testQUeries.GetAccount(context.Background(), accountB.ID)
	require.NoError(t, err)
	fmt.Println("end >>", updatedAccountA.Balance, updatedAccountB.Balance)
}
