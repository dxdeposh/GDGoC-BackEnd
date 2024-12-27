// database.go
package main

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
    dsn := "host=localhost user=postgres password=justher123 dbname=jbeli port=5432 sslmode=disable TimeZone=Asia/Jakarta"
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database!", err)
    }

    // Migrasi schema
    err = DB.AutoMigrate(&User{}, &Category{}, &Product{}, &Order{}, &OrderItem{}, &Transaction{}, &Review{})
    if err != nil {
        log.Fatal("Failed to migrate database!", err)
    }

    fmt.Println("Database connected and migrated successfully.")
}
