package database

import (
	"database/sql"

	"github.com/danielencestari/pos_03/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

// FindAll busca todas as orders com paginação e ordenação
// page: número da página (começando em 1)
// limit: quantidade de registros por página
// sort: campo para ordenação (ex: "id", "price")
func (r *OrderRepository) FindAll(page, limit int, sort string) ([]*entity.Order, error) {
	// Calcula o offset baseado na página
	offset := (page - 1) * limit

	// Monta a query com ordenação e paginação
	// LIMIT controla quantos registros retornar
	// OFFSET controla quantos registros pular
	query := "SELECT id, price, tax, final_price FROM orders ORDER BY " + sort + " LIMIT ? OFFSET ?"

	// Executa a query
	rows, err := r.Db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Importante: sempre fechar o cursor

	// Slice para armazenar as orders encontradas
	var orders []*entity.Order

	// Itera sobre os resultados
	for rows.Next() {
		order := &entity.Order{}
		// Faz o scan dos campos do banco para a struct
		err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return nil, err
		}
		// Adiciona a order ao slice
		orders = append(orders, order)
	}

	// Verifica se houve erro durante a iteração
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetTotal retorna o número total de orders no banco
// Usado para cálculos de paginação
func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("SELECT count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
