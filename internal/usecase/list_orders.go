package usecase

import "github.com/danielencestari/pos_03/internal/entity"

// ListOrdersInputDTO - DTO para receber parâmetros de entrada da listagem
// Segue o mesmo padrão do CreateOrderUseCase
type ListOrdersInputDTO struct {
	Page  int    `json:"page"`  // Número da página (começa em 1)
	Limit int    `json:"limit"` // Quantidade de registros por página
	Sort  string `json:"sort"`  // Campo para ordenação
}

// ListOrdersOutputDTO - DTO para retornar os dados da listagem
// Inclui metadados de paginação além dos dados
type ListOrdersOutputDTO struct {
	Orders     []OrderOutputDTO `json:"orders"`      // Lista de orders
	Page       int              `json:"page"`        // Página atual
	Limit      int              `json:"limit"`       // Limite por página
	Total      int              `json:"total"`       // Total de registros
	TotalPages int              `json:"total_pages"` // Total de páginas
}

// ListOrdersUseCase - Caso de uso para listagem de orders
// Segue o mesmo padrão arquitetural do CreateOrderUseCase
type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface // Interface para acessar dados
}

// NewListOrdersUseCase - Constructor seguindo padrão de injeção de dependência
// Recebe apenas o repository, não precisa de eventos para listagem
func NewListOrdersUseCase(orderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

// Execute - Método principal que executa o caso de uso
// Recebe DTO de entrada e retorna DTO de saída
func (l *ListOrdersUseCase) Execute(input ListOrdersInputDTO) (ListOrdersOutputDTO, error) {
	// Validações de entrada
	if input.Page <= 0 {
		input.Page = 1 // Página padrão
	}
	if input.Limit <= 0 {
		input.Limit = 10 // Limite padrão
	}
	if input.Sort == "" {
		input.Sort = "id" // Ordenação padrão por ID
	}

	// Busca as orders através do repository
	orders, err := l.OrderRepository.FindAll(input.Page, input.Limit, input.Sort)
	if err != nil {
		return ListOrdersOutputDTO{}, err
	}

	// Busca o total de registros para cálculo de paginação
	total, err := l.OrderRepository.GetTotal()
	if err != nil {
		return ListOrdersOutputDTO{}, err
	}

	// Converte as entities Order para DTOs OrderOutput
	var orderDTOs []OrderOutputDTO
	for _, order := range orders {
		dto := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		orderDTOs = append(orderDTOs, dto)
	}

	// Calcula o total de páginas
	totalPages := total / input.Limit
	if total%input.Limit > 0 {
		totalPages++ // Adiciona uma página se houver resto
	}

	// Monta o DTO de saída com todos os dados
	output := ListOrdersOutputDTO{
		Orders:     orderDTOs,
		Page:       input.Page,
		Limit:      input.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return output, nil
}
