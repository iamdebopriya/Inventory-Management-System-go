package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/iamdebopriya/Inventory-Management-System-go/internal/domain"
	"github.com/iamdebopriya/Inventory-Management-System-go/internal/service"
	"github.com/iamdebopriya/Inventory-Management-System-go/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	Usecase *usecase.Usecase
	Email   *service.EmailService
}

func NewHandler(u *usecase.Usecase, e *service.EmailService) *Handler {
	return &Handler{Usecase: u, Email: e}
}

// Root
func (h *Handler) RootRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Inventory API running"})
}

// DTOs

type CreateCategoryInput struct {
	Name string `json:"Name" binding:"required"`
}

type UpdateCategoryInput struct {
	Name string `json:"Name" binding:"required"`
}

type CreateProductInput struct {
	ProductName string  `json:"ProductName" binding:"required"`
	Price       float64 `json:"Price" binding:"required"`
	Description string  `json:"Description"`
	Quantity    int     `json:"Quantity" binding:"required"`
	CategoryID  string  `json:"CategoryID" binding:"required,uuid"`
	IsActive    bool    `json:"IsActive"`
}

type UpdateProductInput struct {
	ProductName string  `json:"ProductName"`
	Price       float64 `json:"Price"`
	Quantity    int     `json:"Quantity"`
	Description string  `json:"Description"`
	CategoryID  string  `json:"CategoryID"`
	IsActive    bool    `json:"IsActive"`
}

type CreateOrderInput struct {
	ProductID string `json:"ProductID" binding:"required,uuid"`
	Quantity  int    `json:"Quantity" binding:"required"`
}

type GetByIDInput struct {
	TableName string `json:"TableName" binding:"required,oneof=categories products orders"`
	ID        string `json:"ID" binding:"required,uuid"`
}

// Categories

func (h *Handler) CreateCategory(c *gin.Context) {
	var input CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := domain.Category{Name: input.Name}
	if err := h.Usecase.CreateCategory(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload from DB to get ID
	cats, _ := h.Usecase.GetCategories()
	for _, cat := range cats {
		if cat.Name == input.Name {
			// Async notification
			go h.Email.CategoryMail(cat.Name)

			c.JSON(http.StatusCreated, cat)
			return
		}
	}

	c.JSON(http.StatusCreated, category)
}

func (h *Handler) GetCategories(c *gin.Context) {
	cats, _ := h.Usecase.GetCategories()
	c.JSON(http.StatusOK, cats)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	var input UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := domain.Category{
		ID:   c.Param("id"),
		Name: input.Name,
	}

	if err := h.Usecase.UpdateCategory(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload to return full entity
	cats, _ := h.Usecase.GetCategories()
	for _, cat := range cats {
		if cat.ID == category.ID {
			c.JSON(http.StatusOK, cat)
			return
		}
	}
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	if err := h.Usecase.DeleteCategory(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Products

func (h *Handler) CreateProduct(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create product entity
	product := domain.Product{
		ProductName: input.ProductName,
		Price:       input.Price,
		Quantity:    input.Quantity,
		CategoryID:  input.CategoryID,
		IsActive:    input.IsActive,
	}

	// Save product in DB (GORM auto-generates ID)
	if err := h.Usecase.Repo.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch the exact product just created, including Category
	var created domain.Product
	if err := h.Usecase.Repo.DB.Preload("Category").First(&created, "id = ?", product.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created product"})
		return
	}

	// Goroutine to simulate notification email
	go h.Email.ProductMail(created.ProductName)

	// Return created product data
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"Product": created,
	})
}

func (h *Handler) GetProducts(c *gin.Context) {
	prods, _ := h.Usecase.GetProducts()
	c.JSON(http.StatusOK, prods)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	var input UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := domain.Product{
		ID:          c.Param("id"),
		ProductName: input.ProductName,
		Price:       input.Price,
		Quantity:    input.Quantity,
		CategoryID:  input.CategoryID,
		IsActive:    input.IsActive,
	}

	if err := h.Usecase.UpdateProduct(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload updated product from DB
	updated, err := h.Usecase.Repo.GetProductByID(product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated product"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	if err := h.Usecase.DeleteProduct(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Orders

func (h *Handler) CreateOrder(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := domain.Order{
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
		OrderDate: time.Time{},
	}

	if err := h.Usecase.CreateOrder(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Reload order from DB to get ID, Product, Category, OrderDate
	orders, _ := h.Usecase.GetOrders()
	var created domain.Order
	for _, o := range orders {
		if o.ProductID == input.ProductID && o.Quantity == input.Quantity {
			created = o
			break
		}
	}

	// Async notification
	go h.Email.OrderMail(created.Product.ProductName, created.Quantity)

	c.JSON(http.StatusCreated, created)
}

func (h *Handler) GetOrders(c *gin.Context) {
	orders, _ := h.Usecase.GetOrders()
	c.JSON(http.StatusOK, orders)
}

// GetTasks returns an aggregated JSON object with categories, products and orders.
// If repository returns a record-not-found error it will be translated to 404.
func (h *Handler) GetTasks(c *gin.Context) {
	tasks, err := h.Usecase.GetTasks()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetByTable is a generic endpoint that fetches a record by ID from any table.
// Expects request body: { "TableName": "categories|products|orders", "ID": "uuid" }
func (h *Handler) GetByTable(c *gin.Context) {
	var input GetByIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch input.TableName {
	case "categories":
		cats, err := h.Usecase.GetCategories()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, cat := range cats {
			if cat.ID == input.ID {
				c.JSON(http.StatusOK, cat)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})

	case "products":
		product, err := h.Usecase.Repo.GetProductByID(input.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, product)

	case "orders":
		orders, err := h.Usecase.GetOrders()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, order := range orders {
			if order.ID == input.ID {
				c.JSON(http.StatusOK, order)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table name"})
	}
}
