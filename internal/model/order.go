package model

import "time"

const (
	OrderStatusCreated   = "created"
	OrderStatusReady     = "ready"
	OrderStatusDone      = "done"
	OrderStatusCancelled = "cancelled"
)

type Order struct {
	ID         int64
	Sum        int
	BuyerID    int64
	ItemsCount int
	Status     string
	CreatedAt  time.Time
}

type OrderExtended struct {
	Order
	Sum   int
	Items []OrderItemInfo
}
