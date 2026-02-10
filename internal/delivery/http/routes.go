package http

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, handler *Handler) {
	r.GET("/", handler.RootRoute)

	api := r.Group("/api/v1")
	{
		// Aggregated tasks endpoint: returns categories, products and orders together
		api.GET("/get-tasks", handler.GetTasks)

		// Generic get-by-tablename endpoint: fetch any record by table name and ID
		api.POST("/get-by-tablename", handler.GetByTable)
		// Categories
		api.GET("/categories", handler.GetCategories)
		api.POST("/categories", handler.CreateCategory)
		api.PUT("/categories/:id", handler.UpdateCategory)
		api.DELETE("/categories/:id", handler.DeleteCategory)

		// Products
		api.GET("/products", handler.GetProducts)
		api.POST("/products", handler.CreateProduct)
		api.PUT("/products/:id", handler.UpdateProduct)
		api.DELETE("/products/:id", handler.DeleteProduct)

		// Orders
		api.GET("/orders", handler.GetOrders)
		api.POST("/orders", handler.CreateOrder)
	}
}
