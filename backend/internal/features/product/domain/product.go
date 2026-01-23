package domain

import (
	"context"
	"dukkanim-api/internal/platform/problem"
	"dukkanim-api/internal/platform/timestamp"
	"dukkanim-api/pkg/viladition"
	"net/http"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID           `json:"id"`
	Slug        string              `json:"slug"`
	Name        string              `json:"name"`
	Sku         string              `json:"sku"` // maybe redundant
	Barcode     string              `json:"barcode"`
	Description string              `json:"description"`
	Category    string              `json:"category"`
	Brand       string              `json:"brand"`
	UnitType    string              `json:"unit_type"`
	IsActive    bool                `json:"is_active"`
	CreatedAt   timestamp.Timestamp `json:"created_at"`
	UpdatedAt   timestamp.Timestamp `json:"updated_at"`
}

type ProductHandlerInterface interface {
	// FindProduct(w http.ResponseWriter, r *http.Request)
	// GetProductList(w http.ResponseWriter, r *http.Request)
	// UpdateProduct(w http.ResponseWriter, r *http.Request)
	// ArchiveProduct(w http.ResponseWriter, r *http.Request)
	CreateProduct(w http.ResponseWriter, r *http.Request)
}

type ProductServiceInterface interface {
	FindProduct(id uuid.UUID, ctx context.Context) (Product, error)
	GetProductList(query ProductFilter, ctx context.Context) ([]Product, error)
	CreateProduct(name, barcode, description, category, brand, unitType string, ctx context.Context) error
	ArchiveProduct(id uuid.UUID, ctx context.Context) error
	UpdateProduct(id uuid.UUID, updated Product, ctx context.Context) error
}

type ProductRepositoryInterface interface {
	FindById(id uuid.UUID, ctx context.Context) (Product, error)
	FetchAll(query ProductFilter, ctx context.Context) ([]Product, error)
	CreateProduct(name, barcode, description, category, brand, unitType string, ctx context.Context) error
	ArchiveProduct(id uuid.UUID, ctx context.Context) error
	UpdateProduct(id uuid.UUID, updated Product, ctx context.Context) error
}

type GetProductListRequest struct {
	Category string `json:"category_id"`
	Sku      string `json:"sku"` // maybe redundant
	Barcode  string `json:"barcode"`
	IsActive bool   `json:"is_active"`
	Name     string `json:"name"`
	Search   string `json:"search"`
	// TODO: Time search
}
type ProductFilter struct {
	// Category string `json:"category_id"`
	// Sku      string `json:"sku"` // maybe redundant
	// Barcode  string `json:"barcode"`
	// IsActive bool   `json:"is_active"`
	// Name     string `json:"name"`
	// Search   string `json:"search"`
	// // TODO: Time search
	GetProductListRequest
}

type GetProductListResponse struct {
	Products []Product
}

type FindProductRequest struct {
	Id      uuid.UUID `json:"product_id"`
	Slug    string    `json:"slug"`
	Barcode string    `json:"barcode"`
	Sku     string    `json:"sku"`
}

type ProductResponse struct {
	Slug        string              `json:"slug"`
	Name        string              `json:"name"`
	Sku         string              `json:"sku"` // maybe redundant
	Barcode     string              `json:"barcode"`
	Description string              `json:"description"`
	Category    string              `json:"category"`
	Brand       string              `json:"brand"`
	UnitType    string              `json:"unit_type"`
	IsActive    bool                `json:"is_active"`
	CreatedAt   timestamp.Timestamp `json:"created_at"`
	UpdatedAt   timestamp.Timestamp `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string `json:"name"`
	Sku         string `json:"sku,omitempty"`
	Barcode     string `json:"barcode"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category"`
	Brand       string `json:"brand"`
	UnitType    string `json:"unit_type"`
}
type UpdateProductRequest struct {
	Id          uuid.UUID `json:"product_id"`
	Name        string    `json:"name"`
	Sku         string    `json:"sku,omitempty"`
	Barcode     string    `json:"barcode"`
	Description string    `json:"description,omitempty"`
	Category    string    `json:"category"`
	Brand       string    `json:"brand"`
	UnitType    string    `json:"unit_type"`
}

type ArchiveProductRequest struct {
	Id uuid.UUID `json:"product_id"`
}

func (r *ProductResponse) FromEntity(product Product) *ProductResponse {
	return &ProductResponse{
		Slug:        product.Slug,
		Name:        product.Name,
		Sku:         product.Sku,
		Barcode:     product.Barcode,
		Description: product.Description,
		Category:    product.Category,
		Brand:       product.Brand,
		UnitType:    product.UnitType,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func (r CreateProductRequest) ToEntity() (*Product, error) {
	var validations map[string]string = make(map[string]string, 0)
	//  Name        string `json:"name"`
	//   Sku         string `json:"sku,omitempty"`
	//   Barcode     string `json:"barcode"`
	//   Description string `json:"description,omitempty"`
	//   Category    string `json:"category"`
	//   Brand       string `json:"brand"`
	//   UnitType    string `json:"unit_type"`

	if err := viladition.
		String(r.Name).
		Required().
		MinLength(5).
		MaxLength(55).
		GetFirstErr(); err != "" {
		validations["name"] = err
	}

	if err := viladition.String(r.Barcode).Required().MinLength(6).MaxLength(50).GetFirstErr(); len(err) != 0 {
		validations["barcode"] = err
	}

	if len(validations) > 0 {
		return nil, problem.UnprocessableEntity.WithValidation(validations)
	}

	return &Product{}, nil
}
