package grpc

import (
	"context"
	"errors"
	"fmt"

	inventoryv1 "github.com/ebisaan/proto/golang/inventory/v1"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ebisaan/inventory/internal/application/core/domain"
)

// GetProductByID implements inventoryv1.InventoryServiceServer.
func (a *Adapter) GetProductByID(ctx context.Context, req *inventoryv1.GetProductByIDRequest) (*inventoryv1.GetProductByIDResponse, error) {
	domainProduct, err := a.app.GetProductByID(ctx, req.Id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return nil, status.New(codes.NotFound, fmt.Sprintf("product with id=%d", req.Id)).Err()
		default:
			zap.L().Error(err.Error())

			return nil, status.New(codes.Unknown, "Unknown Error").Err()
		}
	}

	return &inventoryv1.GetProductByIDResponse{
		Product: protoProduct(domainProduct),
	}, nil
}

// GetProducts implements inventoryv1.InventoryServiceServer.
func (a *Adapter) GetProducts(ctx context.Context, req *inventoryv1.GetProductsRequest) (*inventoryv1.GetProductsResponse, error) {
	domainProducts, meta, err := a.app.GetProducts(ctx, domain.Filter{
		Page:     int(req.GetPagination().GetPage()),
		PageSize: int(req.GetPagination().GetPageSize()),
	})
	if err != nil {
		zap.L().Error(err.Error())

		return nil, status.New(codes.Unknown, "Unknown Error").Err()
	}

	return &inventoryv1.GetProductsResponse{
		Products: protoProducts(domainProducts),
		Metadata: &inventoryv1.Metadata{
			CurrentPage:  int32(meta.CurrentPage),
			FirstPage:    int32(meta.FirstPage),
			LastPage:     int32(meta.LastPage),
			PageSize:     int32(meta.PageSize),
			TotalRecords: int32(meta.TotalRecords),
		},
	}, nil
}

func (a *Adapter) CreateProduct(ctx context.Context, req *inventoryv1.CreateProductRequest) (*inventoryv1.CreateProductResponse, error) {
	id, err := a.app.CreateProduct(ctx, domainProduct(req))
	if err != nil {
		validationErr := &domain.ValidationError{}
		switch {
		case errors.As(err, validationErr):
			return nil, validationErrorToStatusError(validationErr)
		case errors.Is(err, domain.ErrAssociationNotFound):
			st := status.New(codes.InvalidArgument, domain.ErrAssociationNotFound.Error())
			return nil, st.Err()
		default:
			zap.L().Error(err.Error())

			return nil, status.New(codes.Unknown, "Unknown Error").Err()
		}
	}

	return &inventoryv1.CreateProductResponse{
		Id: id,
	}, nil
}

func validationErrorToStatusError(validationErr *domain.ValidationError) error {
	st := status.New(codes.InvalidArgument, "failed validation")
	fieldErrs := make([]*errdetails.BadRequest_FieldViolation, 0, len(validationErr.FieldErrorMessages))
	for field, mess := range validationErr.FieldErrorMessages {
		fieldErrs = append(fieldErrs, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: mess,
		})
	}

	st, _ = st.WithDetails(&errdetails.BadRequest{
		FieldViolations: fieldErrs,
	})

	return st.Err()
}

func protoProducts(dProducts []*domain.Product) []*inventoryv1.Product {
	products := make([]*inventoryv1.Product, 0, len(dProducts))
	for _, dp := range dProducts {
		products = append(products, protoProduct(dp))
	}

	return products
}

func domainProduct(req *inventoryv1.CreateProductRequest) *domain.Product {
	return &domain.Product{
		Name:          req.Name,
		SubCategory:   req.SubCategory,
		StockNumber:   int(req.StockNumber),
		Image:         req.Image,
		DiscountPrice: req.DiscountPrice,
		ActualPrice:   req.ActualPrice,
		CurrencyCode:  req.CurrencyCode,
	}
}

func protoProduct(p *domain.Product) *inventoryv1.Product {
	return &inventoryv1.Product{
		Id:             p.ID,
		Name:           p.Name,
		MainCategory:   p.MainCategory,
		SubCategory:    p.SubCategory,
		StockNumber:    int32(p.StockNumber),
		Image:          p.Image,
		DiscountPrice:  p.DiscountPrice,
		ActualPrice:    p.DiscountPrice,
		CurrencyCode:   p.CurrencyCode,
		CurrencySymbol: p.CurrencySymbol,
	}
}
