package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fvaiiii/ordering_products/order/internal/api/http/handlers"
	"github.com/fvaiiii/ordering_products/order/internal/service"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string, orderService service.OrderService) *Server {
	orderHandler := handlers.NewOrderHandler(orderService)
	router := NewRouter(orderHandler)
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  6 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			println("Server error: ", err.Error())
		}
	}()

	println("Server started on", s.httpServer.Addr)
	<-stop
	println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
