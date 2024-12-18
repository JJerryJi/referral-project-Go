package main

import (
	"fmt"
	"github.com/JerryJi/referral-finder/src/models"
	"github.com/JerryJi/referral-finder/src/routers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitGORM initializes a GORM connection
func InitGORM(dsn string) (*gorm.DB, error) {
	// Open a connection to the MySQL database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the MySQL database successfully with GORM!")
	return db, nil
}

func main() {
	// Data Source Name (DSN)
	dsn := "new_user:password123!@tcp(localhost:3306)/referral?charset=utf8mb4&parseTime=True&loc=Local"

	// Initialize the GORM database connection
	db, err := InitGORM(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//db.AutoMigrate(&User{})

	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	defer rdb.Close()

	db.AutoMigrate(&models.User{}, &models.Student{}, &models.Alumni{})

	// Create a Gin router
	r := gin.Default()

	routers.Register(db, r, rdb)

	// Start the web server
	log.Println("Server running at http://localhost:8080")
	r.Run(":8080") // Default port 8080
}
