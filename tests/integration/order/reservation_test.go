//go:build integration

package tests_order

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"goshop/internal/order/repository"
	productModel "goshop/internal/product/model"
	"goshop/tests/testutil"
)

// TestReserveStock_Concurrency is the headline correctness test for the inventory flow:
// when N goroutines race to reserve from a product with limited stock, exactly `stock`
// reservations must succeed and the rest must fail with ErrInsufficientStock — no oversell.
//
// The implementation relies on a single conditional UPDATE
// ("stock_quantity - reserved_quantity >= ?") which Postgres serializes per-row, so success
// count must equal initial stock regardless of how aggressively the goroutines fight.
func TestReserveStock_Concurrency(t *testing.T) {
	ctx := context.Background()
	db := testutil.StartPostgres(ctx, t)
	require.NoError(t, testutil.ApplyMigrations(db))

	const initialStock = 7
	const racers = 50

	product := &productModel.Product{
		Name:          "race-product",
		Code:          "RACE-1",
		Price:         9.99,
		StockQuantity: initialStock,
		Active:        true,
	}
	require.NoError(t, db.Create(ctx, product))

	repo := repository.NewProductRepository(db)

	var (
		wg         sync.WaitGroup
		ok, failed atomic.Int64
		start      = make(chan struct{})
	)
	for i := 0; i < racers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			err := repo.ReserveStock(ctx, product.ID, 1)
			switch {
			case err == nil:
				ok.Add(1)
			case errors.Is(err, repository.ErrInsufficientStock):
				failed.Add(1)
			default:
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	close(start)
	wg.Wait()

	require.EqualValues(t, initialStock, ok.Load(), "exactly initialStock reservations should win")
	require.EqualValues(t, racers-initialStock, failed.Load(), "remainder should hit ErrInsufficientStock")

	var fresh productModel.Product
	require.NoError(t, db.GetDB().First(&fresh, "id = ?", product.ID).Error)
	require.Equal(t, initialStock, fresh.StockQuantity, "stock_quantity is untouched until commit")
	require.Equal(t, initialStock, fresh.ReservedQuantity, "reserved_quantity matches winners")
}

// TestReserveCommitRelease covers the lifecycle math: reserve 3, commit 1, release 2 →
// stock=initial-1, reserved=0.
func TestReserveCommitRelease(t *testing.T) {
	ctx := context.Background()
	db := testutil.StartPostgres(ctx, t)
	require.NoError(t, testutil.ApplyMigrations(db))

	product := &productModel.Product{
		Name: "lifecycle", Code: "LC-1", Price: 1, StockQuantity: 10, Active: true,
	}
	require.NoError(t, db.Create(ctx, product))

	repo := repository.NewProductRepository(db)
	require.NoError(t, repo.ReserveStock(ctx, product.ID, 3))
	require.NoError(t, repo.CommitReservation(ctx, product.ID, 1))
	require.NoError(t, repo.ReleaseReservation(ctx, product.ID, 2))

	var fresh productModel.Product
	require.NoError(t, db.GetDB().First(&fresh, "id = ?", product.ID).Error)
	require.Equal(t, 9, fresh.StockQuantity)
	require.Equal(t, 0, fresh.ReservedQuantity)
}
