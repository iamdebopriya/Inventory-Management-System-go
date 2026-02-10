package repository

import (
	"errors"
	"inventory-service-go/database"
	"inventory-service-go/internal/domain"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		DB: database.DB,
	}
}

// Category

func (r *Repository) CreateCategory(c domain.Category) error {
	return r.DB.Create(&c).Error
}

func (r *Repository) GetCategories() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.DB.Find(&categories).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return categories, err
}

func (r *Repository) UpdateCategory(c domain.Category) error {
	return r.DB.Model(&domain.Category{}).Where("id = ?", c.ID).Updates(map[string]interface{}{
		"name": c.Name,
	}).Error
}

// Check if any products exist for this category
func (r *Repository) HasProducts(categoryID string) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.Product{}).Where("category_id = ?", categoryID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) DeleteCategory(id string) error {
	hasProducts, err := r.HasProducts(id)
	if err != nil {
		return err
	}
	if hasProducts {
		return errors.New("cannot delete category: products exist")
	}
	return r.DB.Delete(&domain.Category{}, "id = ?", id).Error
}

// Product

func (r *Repository) CreateProduct(p domain.Product) error {
	return r.DB.Create(&p).Error
}

func (r *Repository) GetProducts() ([]domain.Product, error) {
	var products []domain.Product
	err := r.DB.Preload("Category").Find(&products).Error
	return products, err
}

func (r *Repository) GetProductByID(id string) (*domain.Product, error) {
	var p domain.Product
	err := r.DB.Preload("Category").First(&p, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &p, err
}

func (r *Repository) UpdateProduct(p domain.Product) error {
	return r.DB.Model(&domain.Product{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
		"product_name": p.ProductName,
		"price":        p.Price,
		"description":  p.Description,
		"quantity":     p.Quantity,
		"category_id":  p.CategoryID,
		"is_active":    p.IsActive,
	}).Error
}

func (r *Repository) DeleteProduct(id string) error {
	var count int64
	r.DB.Model(&domain.Order{}).Where("product_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("cannot delete product: orders exist")
	}
	return r.DB.Delete(&domain.Product{}, "id = ?", id).Error
}

// Order

func (r *Repository) CreateOrder(o domain.Order) error {
	return r.DB.Create(&o).Error
}

func (r *Repository) GetOrders() ([]domain.Order, error) {
	var orders []domain.Order
	err := r.DB.Preload("Product").Preload("Product.Category").Find(&orders).Error
	return orders, err
}
