package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Bioskop struct {
	ID     int
	Nama   string
	Lokasi string
	Rating float32
}

// Endpoint POST
func CreateBioskop(c *gin.Context, db *sql.DB) {
	BasicAuth()(c)
	if c.IsAborted() {
		return
	}
	var newBioskop Bioskop
	if err := c.ShouldBindJSON(&newBioskop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	if newBioskop.Nama == "" || newBioskop.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Lokasi tidak boleh kosong"})
		return
	}
	query := `INSERT INTO bioskop ("Nama", "Lokasi", "Rating") VALUES ($1, $2, $3) RETURNING "ID"`
	err := db.QueryRow(query, newBioskop.Nama, newBioskop.Lokasi, newBioskop.Rating).Scan(&newBioskop.ID)
	if err != nil {
		log.Println("Error inserting new bioskop:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan bioskop"})
		return
	}
	c.JSON(http.StatusCreated, newBioskop)
}

// Endpoint GET
func GetBioskops(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`SELECT "ID", "Nama", "Lokasi", "Rating" FROM bioskop`)
	if err != nil {
		log.Println("Error fetching bioskop:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data bioskop"})
		return
	}
	defer rows.Close()

	var bioskops []Bioskop
	for rows.Next() {
		var b Bioskop
		if err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating); err != nil {
			log.Println("Error scanning bioskop:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data bioskop"})
			return
		}
		bioskops = append(bioskops, b)
	}
	c.JSON(http.StatusOK, bioskops)
}

// Endpoint GET by Id
func GetBioskopByID(c *gin.Context, db *sql.DB) {
	id := c.Param("ID")
	var b Bioskop
	query := `SELECT "ID", "Nama", "Lokasi", "Rating" FROM bioskop WHERE "ID" = $1`
	err := db.QueryRow(query, id).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bioskop tidak ditemukan"})
		} else {
			log.Println("Error fetching bioskop by ID:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data bioskop"})
		}
		return
	}
	c.JSON(http.StatusOK, b)
}

// basic auth middleware
func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, password, hasAuth := c.Request.BasicAuth()
		if hasAuth && user == "admin" && password == "root" {
			c.Next()
			return
		}
		c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}