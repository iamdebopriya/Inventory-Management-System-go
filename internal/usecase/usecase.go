package usecase

import (
	"errors"
	"inventory-service-go/internal/domain"
	"inventory-service-go/internal/repository"
)

type Usecase struct {
	Repo *repository.Repository
}

func NewUsecase(r *repository.Repository) *Usecase {
	return &Usecase{Repo: r}
}

// Category

func (u *Usecase) CreateCategory(c domain.Category) error {
	return u.Repo.CreateCategory(c)
}

func (u *Usecase) GetCategories() ([]domain.Category, error) {
	return u.Repo.GetCategories()
}

func (u *Usecase) UpdateCategory(c domain.Category) error {
	return u.Repo.UpdateCategory(c)
}

func (u *Usecase) DeleteCategory(id string) error {
	return u.Repo.DeleteCategory(id)
}

// Product

func (u *Usecase) CreateProduct(p domain.Product) error {
	return u.Repo.CreateProduct(p)
}

func (u *Usecase) GetProducts() ([]domain.Product, error) {
	return u.Repo.GetProducts()
}

func (u *Usecase) UpdateProduct(p domain.Product) error {
	return u.Repo.UpdateProduct(p)
}

func (u *Usecase) DeleteProduct(id string) error {
	return u.Repo.DeleteProduct(id)
}

// Order

func (u *Usecase) CreateOrder(o domain.Order) error {
	product, err := u.Repo.GetProductByID(o.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Quantity < o.Quantity {
		return errors.New("product quantity insufficient")
	}

	product.Quantity -= o.Quantity
	if err := u.Repo.UpdateProduct(*product); err != nil {
		return err
	}

	return u.Repo.CreateOrder(o)
}

func (u *Usecase) GetOrders() ([]domain.Order, error) {
	return u.Repo.GetOrders()
}

type TasksResponse struct {
	Categories []domain.Category `json:"categories"`
	Products   []domain.Product  `json:"products"`
	Orders     []domain.Order    `json:"orders"`
}

// GetTasks aggregates categories, products and orders into a single response.
func (u *Usecase) GetTasks() (*TasksResponse, error) {
	// Reuse repository methods which already encapsulate DB access.
	cats, err := u.Repo.GetCategories()
	if err != nil {
		return nil, err
	}

	prods, err := u.Repo.GetProducts()
	if err != nil {
		return nil, err
	}

	orders, err := u.Repo.GetOrders()
	if err != nil {
		return nil, err
	}

	return &TasksResponse{
		Categories: cats,
		Products:   prods,
		Orders:     orders,
	}, nil
}
