package web

import (
	"encoding/json"
	"net/http"
	"strconv" // Para converter strings em números

	"github.com/danielencestari/pos_03/internal/entity"
	"github.com/danielencestari/pos_03/internal/usecase"
	"github.com/danielencestari/pos_03/pkg/events"
)

type WebOrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
	}
}

// Create - Handler para criar uma nova order via POST /order
// Segue o padrão REST: POST para criação
func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Define o Content-Type da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// 1. Decodifica o JSON do request body para DTO
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Cria o use case com suas dependências
	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)

	// 3. Executa o use case
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Retorna o resultado como JSON
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// List - Handler para listar orders via GET /order
// Segue o padrão REST: GET para leitura
func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	// Define o Content-Type da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// 1. Extrai parâmetros de query da URL
	// Ex: GET /order?page=1&limit=10&sort=id
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	// 2. Converte strings para números com valores padrão
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if sort == "" {
		sort = "id" // Ordenação padrão
	}

	// 3. Monta o DTO de entrada
	dto := usecase.ListOrdersInputDTO{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	// 4. Cria e executa o use case
	listOrders := usecase.NewListOrdersUseCase(h.OrderRepository)
	output, err := listOrders.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Retorna o resultado como JSON
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
