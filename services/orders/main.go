package orders

import (
	"golang-grpc/internal/server"
	"golang-grpc/services/orders/store"
)

type OrdersService struct {
	config *store.InitialConfig
}

func NewOrdersService(config *store.InitialConfig) *OrdersService {
	return &OrdersService{
		config: config,
	}
}

func (os *OrdersService) ExecuteExternal() {
	os.config = &store.InitialConfig{}

	ready := make(chan bool, 1)
	done := os.Execute(ready)
	<-done
}

func (os *OrdersService) Execute(ready chan<- bool) <-chan bool {
	var grpcServer server.Server = NewGRPCServer(&gRPCServerConfig{
		ServerConfig: server.ServerConfig{
			Host: "localhost",
			Port: 3081,
		},
	})
	var httpServer server.Server = NewHTTPServer(&httpServerConfig{
		ServerConfig: server.ServerConfig{
			Host: "localhost",
			Port: 3082,
		},
	})

	runList := []server.Server{grpcServer, httpServer}
	server.RunServersInParallel(runList, server.ServerRunConfig{
		WithGracefulShutdown: true,
		Silent:               true,
	})

	ready <- true
	doneChannel := make(chan bool, 1)
	go func() {
		select {
		case doneChannel <- <-grpcServer.GetDoneChannel():
			break
		case doneChannel <- <-httpServer.GetDoneChannel():
			break
		}
	}()

	return doneChannel
}
