// main.go
package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// main.go (modifikasi)
func main() {
	InitDatabase()

	router := gin.Default()
	router.Use(cors.Default())

	// Public endpoints
	router.POST("/users", CreateUser)
	router.POST("/login", Login)

	// Protected routes
	protected := router.Group("/", AuthMiddleware())
	{
			// User endpoints
			protected.GET("/users", GetUsers)
			protected.GET("/users/:id", GetUser)
			protected.PUT("/users/:id", UpdateUser)
			protected.DELETE("/users/:id", DeleteUser)

			// Category endpoints
			protected.POST("/categories", CreateCategory)
			protected.GET("/categories", GetCategories)
			protected.GET("/categories/:id", GetCategory)
			protected.PUT("/categories/:id", UpdateCategory)
			protected.DELETE("/categories/:id", DeleteCategory)

			// Product endpoints
			protected.POST("/products", CreateProduct)
			protected.GET("/products", GetProduct)
			protected.GET("/products/:id", GetProduct)
			protected.PUT("/products/:id", UpdateProduct)
			protected.DELETE("/products/:id", DeleteProduct)

			// Order endpoints
			protected.POST("/orders", CreateOrder)
			protected.GET("/orders", GetOrders)
			protected.GET("/orders/:id", GetOrder)
			protected.PUT("/orders/:id", UpdateOrder)
			protected.DELETE("/orders/:id", DeleteOrder)

			// Review endpoints
			protected.POST("/reviews", CreateReview)
			protected.GET("/reviews", GetReviews)
			protected.GET("/reviews/:id", GetReview)
			protected.PUT("/reviews/:id", UpdateReview)
			protected.DELETE("/reviews/:id", DeleteReview)

			// Transaction endpoints
			protected.POST("/transactions", CreateTransaction)
			protected.GET("/transactions", GetTransactions)
			protected.GET("/transactions/:id", GetTransaction)
			protected.PUT("/transactions/:id", UpdateTransaction)
			protected.DELETE("/transactions/:id", DeleteTransaction)
	}

	router.Run(":8080")
}

