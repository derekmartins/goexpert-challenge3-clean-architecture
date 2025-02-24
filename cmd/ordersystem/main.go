package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/configs"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/event/handler"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/infra/graph"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/infra/grpc/pb"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/infra/grpc/service"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/internal/infra/web/webserver"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	appConfigs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(appConfigs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", appConfigs.DBUser, appConfigs.DBPassword, appConfigs.DBHost, appConfigs.DBPort, appConfigs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()

	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersOutputUseCase := NewListOrdersOutputUseCase(db)

	appWebserver := webserver.NewWebServer(appConfigs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	appWebserver.AddHandler("/order", webOrderHandler.Create)
	fmt.Println("Starting web server on port", appConfigs.WebServerPort)
	go appWebserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", appConfigs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfigs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase:      *createOrderUseCase,
		ListOrdersOutputUseCase: *listOrdersOutputUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", appConfigs.GraphQLServerPort)
	http.ListenAndServe(":"+appConfigs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
