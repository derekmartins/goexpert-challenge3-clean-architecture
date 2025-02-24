package usecase

import (
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/entity"
)

type OrderOutputListDTO []OrderOutputDTO

type ListOrdersOutputUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersOutputUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrdersOutputUseCase {
	return &ListOrdersOutputUseCase{
		OrderRepository: OrderRepository,
	}
}

func (listOrdersOutputUseCase *ListOrdersOutputUseCase) Execute() (OrderOutputListDTO, error) {
	listOrderOutput, err := listOrdersOutputUseCase.OrderRepository.List()
	if err != nil {
		return OrderOutputListDTO{}, err
	}

	var listOrdersOutputDTO OrderOutputListDTO

	for _, order := range listOrderOutput {
		listOrdersOutputDTO = append(listOrdersOutputDTO, OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return listOrdersOutputDTO, nil
}
