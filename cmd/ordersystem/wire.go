//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/entity"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/event"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/infra/database"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/infra/web"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/usecase"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrdersOutputUseCase(db *sql.DB) *usecase.ListOrdersOutputUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrdersOutputUseCase,
	)
	return &usecase.ListOrdersOutputUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
