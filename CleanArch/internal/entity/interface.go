package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	FindAll(page, limit int, sort string) ([]*Order, error) // Nova função para listagem com paginação
	GetTotal() (int, error)                                 // Descomentada para contar total de registros
}
