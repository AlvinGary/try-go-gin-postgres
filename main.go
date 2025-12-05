package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"try-go-gin-postgres/database"
	"try-go-gin-postgres/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	err error
)

func main() {
	err = godotenv.Load("config/.env")
    if err != nil {
        panic("Error loading .env file")
    }
	dbInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot reach the database:", err)
	}

	database.DBMigrate(db)
	router := gin.Default()
	routers.SetupBioskopRoutes(router, db)
	log.Println("Server is running on port 8000")
	router.Run(":8000")
}