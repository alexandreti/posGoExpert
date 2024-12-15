package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	ListOrders() ([]Order, error)
}
