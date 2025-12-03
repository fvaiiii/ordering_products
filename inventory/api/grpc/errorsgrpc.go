package grpc

import (
	"errors"

	"github.com/fvaiiii/ordering_products/inventory/internal/port/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	if errors.Is(err, repo.ErrNotFound) {
		return status.Error(codes.NotFound, "product is not found")
	}
	return status.Error(codes.Internal, "internal server error")
}
