package main

import (
	"database/sql"
	"log"
	"try-go-gin-postgres/routers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "root"
	DB_NAME     = "try-postgres"
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
)

var (
	db *sql.DB
	err error
)

func main() {
	dbInfo := "host=" + DB_HOST + " port=" + DB_PORT + " user=" + DB_USER + " password=" + DB_PASSWORD + " dbname=" + DB_NAME + " sslmode=disable"
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot reach the database:", err)
	}
	router := gin.Default()
	routers.SetupBioskopRoutes(router, db)
	log.Println("Server is running on port 8000")
	router.Run(":8000")
}