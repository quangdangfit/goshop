package service

import "fmt"

// InsufficientStockError carries the offending product so the HTTP handler can return a
// structured 409 body the FE uses to reconcile the local cart.
type InsufficientStockError struct {
	ProductID string
	Requested int
}

func (e *InsufficientStockError) Error() string {
	return fmt.Sprintf("insufficient stock for product %s (requested %d)", e.ProductID, e.Requested)
}
