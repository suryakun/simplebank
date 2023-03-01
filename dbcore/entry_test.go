package dbcore

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createMockEntry(t *testing.T) (Entry, Account) {
	account := createMockAccount(t)
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    10,
	}

	entry, err := testQUeries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotEmpty(t, entry.ID)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	require.NotEmpty(t, entry.AccountID)
	require.NotEmpty(t, entry.Amount)
	return entry, account
}
func TestCreateEntry(t *testing.T) {
	createMockEntry(t)
}

func TestGetEntry(t *testing.T) {
	entryA, _ := createMockEntry(t)
	entryB, err := testQUeries.GetEntry(context.Background(), entryA.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryB)

	require.Equal(t, entryA.ID, entryB.ID)
	require.Equal(t, entryA.AccountID, entryB.AccountID)
	require.Equal(t, entryA.Amount, entryB.Amount)
	require.WithinDuration(t, entryA.CreatedAt, entryB.CreatedAt, time.Second)
}

func TestUpdateEntryy(t *testing.T) {
	entryA, _ := createMockEntry(t)
	args := UpdateEntryParams{
		ID:     entryA.ID,
		Amount: 20,
	}

	entry, err := testQUeries.UpdateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.Amount, args.Amount)
}

func TestDeleteEntry(t *testing.T) {
	entry, _ := createMockEntry(t)
	err := testQUeries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	entryA, err := testQUeries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.Empty(t, entryA)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		_, _ = createMockEntry(t)
	}
	args := GetEntriesParams{
		Offset: 5,
		Limit:  5,
	}

	entries, err := testQUeries.GetEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
