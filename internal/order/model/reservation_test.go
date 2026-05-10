package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStockReservationBeforeCreateSetsDefaults(t *testing.T) {
	r := &StockReservation{}
	require.NoError(t, r.BeforeCreate(nil))
	require.NotEmpty(t, r.ID)
	require.Equal(t, ReservationStatusActive, r.Status)
}

func TestStockReservationBeforeCreatePreservesExisting(t *testing.T) {
	r := &StockReservation{ID: "fixed", Status: ReservationStatusCommitted}
	require.NoError(t, r.BeforeCreate(nil))
	require.Equal(t, "fixed", r.ID)
	require.Equal(t, ReservationStatusCommitted, r.Status)
}
