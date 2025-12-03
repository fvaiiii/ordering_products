package grpc

import (
	"context"

	"github.com/fvaiiii/ordering_products/inventory/internal/service"
	v1 "github.com/fvaiiii/ordering_products/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	v1.UnimplementedInventoryServiceServer
	service *service.InventoryService
}

func NewServer(svc *service.InventoryService) *Server {
	return &Server{
		service: svc,
	}
}

func (s *Server) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.GetProductResponse, error) {
	if req.GetUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	product, err := s.service.GetProduct(ctx, req.GetUuid())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &v1.GetProductResponse{
		Product: domainToProto(product),
	}, nil
}

func (s *Server) ListProducts(ctx context.Context, req *v1.ListProductsRequest) (*v1.ListProductsResponse, error) {
	products, err := s.service.ListProducts(ctx, filterFromProto(req.GetFilter()))
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &v1.ListProductsResponse{
		Products: domainArrToProto(products),
	}, nil

}
