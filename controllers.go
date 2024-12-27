// controllers.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// =====================
// User Controllers
// =====================

// CreateUser
func CreateUser(c *gin.Context) {
    var input User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    input.Password = string(hashedPassword)

    if err := DB.Create(&input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Jangan kembalikan password
    input.Password = ""
    c.JSON(http.StatusOK, gin.H{"data": input})
}

// GetUsers
func GetUsers(c *gin.Context) {
    var users []User
    DB.Find(&users)
    // Jangan kembalikan password
    for i := range users {
        users[i].Password = ""
    }
    c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetUser
func GetUser(c *gin.Context) {
    var user User
    if err := DB.First(&user, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    // Jangan kembalikan password
    user.Password = ""
    c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser
func UpdateUser(c *gin.Context) {
    var user User
    if err := DB.First(&user, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    var input User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Jika password diupdate, hash terlebih dahulu
    if input.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
            return
        }
        input.Password = string(hashedPassword)
    } else {
        input.Password = user.Password // Pertahankan password lama
    }

    DB.Model(&user).Updates(input)
    // Jangan kembalikan password
    user.Password = ""
    c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser
func DeleteUser(c *gin.Context) {
    var user User
    if err := DB.First(&user, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    DB.Delete(&user)
    c.JSON(http.StatusOK, gin.H{"data": true})
}

// =====================
// Category Controllers
// =====================

// CreateCategory
func CreateCategory(c *gin.Context) {
    var input Category
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := DB.Create(&input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": input})
}

// GetCategories
func GetCategories(c *gin.Context) {
    var categories []Category
    DB.Find(&categories)
    c.JSON(http.StatusOK, gin.H{"data": categories})
}

// GetCategory
func GetCategory(c *gin.Context) {
    var category Category
    if err := DB.First(&category, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": category})
}

// UpdateCategory
func UpdateCategory(c *gin.Context) {
    var category Category
    if err := DB.First(&category, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    var input Category
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    DB.Model(&category).Updates(input)
    c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteCategory
func DeleteCategory(c *gin.Context) {
    var category Category
    if err := DB.First(&category, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    DB.Delete(&category)
    c.JSON(http.StatusOK, gin.H{"data": true})
}

// =====================
// Product Controllers
// =====================

// CreateProduct
func CreateProduct(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
	}
	user := userInterface.(User)

	if user.Role != "seller" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only sellers can create products"})
			return
	}

	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	input.UserID = user.ID

	if err := DB.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	c.JSON(http.StatusOK, gin.H{"data": input})
}

// GetProducts
func GetProducts(c *gin.Context) {
    var products []Product
    DB.Preload("Category").Preload("User").Find(&products)
    c.JSON(http.StatusOK, gin.H{"data": products})
}

// GetProduct
func GetProduct(c *gin.Context) {
    var product Product
    if err := DB.Preload("Category").Preload("User").First(&product, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": product})
}

// UpdateProduct
func UpdateProduct(c *gin.Context) {
    var product Product
    if err := DB.First(&product, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    user := userInterface.(User)

    if product.UserID != user.ID && user.Role != "admin" { // Asumsi ada role "admin"
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this product"})
        return
    }

    var input Product
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    DB.Model(&product).Updates(input)
    c.JSON(http.StatusOK, gin.H{"data": product})
}

// DeleteProduct
func DeleteProduct(c *gin.Context) {
    var product Product
    if err := DB.First(&product, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    user := userInterface.(User)

    if product.UserID != user.ID && user.Role != "admin" { // Asumsi ada role "admin"
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this product"})
        return
    }

    DB.Delete(&product)
    c.JSON(http.StatusOK, gin.H{"data": true})
}

// =====================
// Order Controllers
// =====================

// CreateOrder
func CreateOrder(c *gin.Context) {
    var input struct {
        OrderItems []struct {
            ProductID uint `json:"product_id" binding:"required"`
            Quantity  int  `json:"quantity" binding:"required,gt=0"`
        } `json:"order_items" binding:"required,dive,required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    user := userInterface.(User)

    var totalPrice float64
    var orderItems []OrderItem

    for _, item := range input.OrderItems {
        var product Product
        if err := DB.First(&product, item.ProductID).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
            return
        }

        if product.Stock < item.Quantity {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product ID " + string(product.ID)})
            return
        }

        totalPrice += product.Price * float64(item.Quantity)
        orderItems = append(orderItems, OrderItem{
            ProductID: product.ID,
            Quantity:  item.Quantity,
            Price:     product.Price,
        })
    }

    order := Order{
        UserID:     user.ID,
        TotalPrice: totalPrice,
        Status:     "pending",
        OrderItems: orderItems,
    }

    if err := DB.Create(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Kurangi stok produk
    for _, item := range orderItems {
        DB.Model(&Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock - ?", item.Quantity))
    }

    c.JSON(http.StatusOK, gin.H{"data": order})
}

// GetOrders
func GetOrders(c *gin.Context) {
    var orders []Order
    DB.Preload("OrderItems.Product").Preload("Transaction").Preload("User").Find(&orders)
    c.JSON(http.StatusOK, gin.H{"data": orders})
}

// GetOrder
func GetOrder(c *gin.Context) {
    var order Order
    if err := DB.Preload("OrderItems.Product").Preload("Transaction").Preload("User").First(&order, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": order})
}

// UpdateOrder
func UpdateOrder(c *gin.Context) {
    var order Order
    if err := DB.First(&order, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    var input Order
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Contoh: hanya status yang bisa diupdate
    DB.Model(&order).Updates(map[string]interface{}{"status": input.Status})
    c.JSON(http.StatusOK, gin.H{"data": order})
}

// DeleteOrder
func DeleteOrder(c *gin.Context) {
    var order Order
    if err := DB.First(&order, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    DB.Delete(&order)
    c.JSON(http.StatusOK, gin.H{"data": true})
}

// =====================
// Review Controllers
// =====================

// CreateReview
func CreateReview(c *gin.Context) {
    var input Review
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    user := userInterface.(User)

    input.UserID = user.ID

    if err := DB.Create(&input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": input})
}

// GetReviews
func GetReviews(c *gin.Context) {
    var reviews []Review
    DB.Preload("User").Preload("Product").Find(&reviews)
    c.JSON(http.StatusOK, gin.H{"data": reviews})
}

// GetReview
func GetReview(c *gin.Context) {
    var review Review
    if err := DB.Preload("User").Preload("Product").First(&review, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": review})
}

// UpdateReview
func UpdateReview(c *gin.Context) {
    var review Review
    if err := DB.First(&review, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    user := userInterface.(User)

    if review.UserID != user.ID && user.Role != "admin" { // Asumsi ada role "admin"
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this review"})
        return
    }

    var input Review
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    DB.Model(&review).Updates(input)
    c.JSON(http.StatusOK, gin.H{"data": review})
}

// DeleteReview
func DeleteReview(c *gin.Context) {
    var review Review
    if err := DB.First(&review, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    user := userInterface.(User)

    if review.UserID != user.ID && user.Role != "admin" { // Asumsi ada role "admin"
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this review"})
        return
    }

    DB.Delete(&review)
    c.JSON(http.StatusOK, gin.H{"data": true})
}

// =====================
// Transaction Controllers
// =====================

// CreateTransaction
func CreateTransaction(c *gin.Context) {
    var input Transaction
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := DB.Create(&input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": input})
}

// GetTransactions
func GetTransactions(c *gin.Context) {
    var transactions []Transaction
    DB.Preload("Order").Find(&transactions)
    c.JSON(http.StatusOK, gin.H{"data": transactions})
}

// GetTransaction
func GetTransaction(c *gin.Context) {
    var transaction Transaction
    if err := DB.Preload("Order").First(&transaction, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// UpdateTransaction
func UpdateTransaction(c *gin.Context) {
    var transaction Transaction
    if err := DB.First(&transaction, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    var input Transaction
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    DB.Model(&transaction).Updates(input)
    c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// DeleteTransaction
func DeleteTransaction(c *gin.Context) {
    var transaction Transaction
    if err := DB.First(&transaction, c.Param("id")).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    DB.Delete(&transaction)
    c.JSON(http.StatusOK, gin.H{"data": true})
}
