// models.go
package main

import (
	"time"
)

type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `json:"name" binding:"required"`
    Email     string         `gorm:"unique" json:"email" binding:"required,email"`
    Password  string         `json:"-"`
    Role      string         `json:"role"` // e.g., "buyer", "seller"
    Products  []Product      `json:"products"`
    Orders    []Order        `json:"orders"`
    Reviews   []Review       `json:"reviews"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
}

type Category struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `json:"name" binding:"required"`
    Description string    `json:"description"`
    Products    []Product `json:"products"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `json:"name" binding:"required"`
    Description string    `json:"description"`
    Price       float64   `json:"price" binding:"required,gt=0"`
    Stock       int       `json:"stock" binding:"required,gt=0"`
    CategoryID  uint      `json:"category_id" binding:"required"`
    Category    Category  `json:"category"`
    UserID      uint      `json:"user_id"`
    User        User      `json:"user"`
    Reviews     []Review  `json:"reviews"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Order struct {
    ID          uint         `gorm:"primaryKey" json:"id"`
    UserID      uint         `json:"user_id"`
    User        User         `json:"user"`
    TotalPrice  float64      `json:"total_price"`
    Status      string       `json:"status"` // e.g., "pending", "completed", "cancelled"
    OrderItems  []OrderItem  `json:"order_items"`
    Transaction *Transaction `json:"transaction"` // Menggunakan pointer untuk menghindari rekursi
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
}

type OrderItem struct {
    ID        uint     `gorm:"primaryKey" json:"id"`
    OrderID   uint     `json:"order_id"`
    Order     *Order   `json:"order"` // Menggunakan pointer untuk menghindari rekursi
    ProductID uint     `json:"product_id"`
    Product   Product  `json:"product"`
    Quantity  int      `json:"quantity" binding:"required,gt=0"`
    Price     float64  `json:"price"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    OrderID         uint      `json:"order_id"`
    Order           *Order    `json:"order"` // Menggunakan pointer untuk menghindari rekursi
    PaymentMethod   string    `json:"payment_method"` // e.g., "credit_card", "paypal"
    PaymentStatus   string    `json:"payment_status"` // e.g., "paid", "unpaid"
    TransactionDate time.Time `json:"transaction_date"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type Review struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `json:"user_id"`
    User      User      `json:"user"`
    ProductID uint      `json:"product_id"`
    Product   Product   `json:"product"`
    Rating    int       `json:"rating" binding:"required,gte=1,lte=5"`
    Comment   string    `json:"comment"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
