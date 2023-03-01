package dbcore

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/suryakun/simplebank/util"
)

func createMockAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  500,
		Currency: util.RandomCurrency(),
	}

	account, err := testQUeries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.NotEmpty(t, account.ID)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.NotEmpty(t, account.Currency)
	require.NotEmpty(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	createMockAccount(t)
}

func TestGetAccount(t *testing.T) {
	accountA := createMockAccount(t)
	accountB, err := testQUeries.GetAccount(context.Background(), accountA.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountB)

	require.Equal(t, accountA.ID, accountB.ID)
	require.Equal(t, accountA.Owner, accountB.Owner)
	require.Equal(t, accountA.Balance, accountB.Balance)
	require.Equal(t, accountA.Currency, accountB.Currency)
	require.WithinDuration(t, accountA.CreatedAt, accountB.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	accountA := createMockAccount(t)
	args := UpdateAccountParams{ID: accountA.ID, Balance: util.RandomMoney()}

	accountB, err := testQUeries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accountB)
	require.Equal(t, accountA.ID, accountB.ID)
	require.Equal(t, accountB.Owner, accountB.Owner)
	require.Equal(t, accountB.Currency, accountA.Currency)
	require.Equal(t, accountB.Balance, args.Balance)
	require.WithinDuration(t, accountA.CreatedAt, accountB.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	accountA := createMockAccount(t)
	err := testQUeries.DeleteAccount(context.Background(), accountA.ID)
	require.NoError(t, err)
	accountB, err := testQUeries.GetAccount(context.Background(), accountA.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountB)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createMockAccount(t)
	}

	args := GetAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQUeries.GetAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
