package constants

const (
	LimitDataPerPage = 10
	LimitInsertBatch = 100
)

type OrderStatus uint

const (
	OrderStatusWaitingPayment OrderStatus = iota + 1
	OrderStatusPending
	OrderStatusProcessing
	OrderStatusShipped
	OrderStatusDelivered
	OrderStatusCancelled
	OrderStatusRefunded
	OrderStatusFailed
)

func (o OrderStatus) String() string {
	return [...]string{"", "WAITINGPAYMENT", "PENDING", "PROCESSING", "SHIPPED", "DELIVERED", "CANCELLED", "REFUNDED", "FAILED"}[o]
}

func (o OrderStatus) EnumIndex() uint {
	return uint(o)
}
