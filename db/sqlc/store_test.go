package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	const n = 5
	const amount int64 = 10
	errs := make(chan error, n)
	results := make(chan *TransferTxResult, n)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	executed := make(map[int]bool)
	for i := 0; i < n; i++ {
		require.NoError(t, <-errs)
		result := <-results
		require.NotNil(t, result)

		transfer := result.Transfer
		require.NotZero(t, transfer.ID)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		dbTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		require.Equal(t, transfer, dbTransfer)

		fromEntry := result.FromEntry
		require.NotZero(t, fromEntry.ID)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		toEntry := result.ToEntry
		require.NotZero(t, toEntry.ID)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		fromDiff := fromAccount.Balance - result.FromAccount.Balance
		toDiff := result.ToAccount.Balance - toAccount.Balance
		require.Equal(t, fromDiff, toDiff)
		require.Positive(t, fromDiff)
		require.Zero(t, fromDiff%amount)

		step := int(fromDiff / amount)
		require.GreaterOrEqual(t, step, 1)
		require.LessOrEqual(t, step, n)
		require.NotContains(t, executed, step)
		executed[step] = true
	}

	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedToAccount.Balance)
}
