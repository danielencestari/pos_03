package service

import (
	"context"

	"github.com/danielencestari/pos_03/internal/infra/grpc/pb"
	"github.com/danielencestari/pos_03/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase // Adicionando o novo use case
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrdersUseCase usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

// CreateOrder - Método gRPC para criar uma order
// Converte proto request → DTO → use case → DTO → proto response
func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// 1. Converte proto request para DTO
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}

	// 2. Executa o use case
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	// 3. Converte DTO para proto response
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

// ListOrders - Método gRPC para listar orders com paginação
// Segue o mesmo padrão: proto → DTO → use case → DTO → proto
func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	// 1. Converte proto request para DTO de entrada
	dto := usecase.ListOrdersInputDTO{
		Page:  int(in.Page),  // Converte int32 para int
		Limit: int(in.Limit), // Converte int32 para int
		Sort:  in.Sort,
	}

	// 2. Executa o use case
	output, err := s.ListOrdersUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	// 3. Converte slice de OrderOutputDTO para slice de OrderItem (proto)
	var orderItems []*pb.OrderItem
	for _, order := range output.Orders {
		item := &pb.OrderItem{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		}
		orderItems = append(orderItems, item)
	}

	// 4. Monta a resposta proto com os dados de paginação
	return &pb.ListOrdersResponse{
		Orders:     orderItems,               // Lista de orders
		Page:       int32(output.Page),       // Converte int para int32
		Limit:      int32(output.Limit),      // Converte int para int32
		Total:      int32(output.Total),      // Total de registros
		TotalPages: int32(output.TotalPages), // Total de páginas
	}, nil
}
