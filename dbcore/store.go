package dbcore

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	transaction, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	queries := New(transaction)

	err = fn(queries)
	if err != nil {
		if errRollback := transaction.Rollback(); errRollback != nil {
			return fmt.Errorf("error rollback %v, with error : %v", err, errRollback)
		}
		return err
	}
	return transaction.Commit()
}

type TransferTransactionParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTransactionResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTransaction(ctx context.Context, args TransferTransactionParams) (TransferTransactionResult, error) {
	var result TransferTransactionResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(args))
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = makeTransactionBalance(ctx, q, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount)
		} else {
			result.FromAccount, result.ToAccount, err = makeTransactionBalance(ctx, q, args.ToAccountID, -args.Amount, args.FromAccountID, args.Amount)
		}

		return err
	})

	return result, err
}

func makeTransactionBalance(ctx context.Context, q *Queries, fromAccountID int64, fromAmount int64, toAccountID int64, toAmount int64) (fromAccount Account, toAccount Account, err error) {
	fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccountID,
		Amount: fromAmount,
	})
	if err != nil {
		return
	}

	toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccountID,
		Amount: toAmount,
	})
	if err != nil {
		return
	}
	return
}
