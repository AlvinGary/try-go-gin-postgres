package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"try-go-gin-postgres/structs"

	"github.com/gin-gonic/gin"
)

// Endpoint POST
func CreateBioskop(c *gin.Context, db *sql.DB) {
	BasicAuth()(c)
	if c.IsAborted() {
		return
	}
	var newBioskop structs.Bioskop
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

	var bioskops []structs.Bioskop
	for rows.Next() {
		var b structs.Bioskop
		if err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating); err != nil {
			log.Println("Error scanning bioskop:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data bioskop"})
			return
		}
		bioskops = append(bioskops, b)
	}
	c.JSON(http.StatusOK, bioskops)
}

// Endpoint GET by ID
func GetBioskopByID(c *gin.Context, db *sql.DB) {
	id := c.Param("ID")
	var b structs.Bioskop
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

// Endpoint PUT by ID
func UpdateBioskop(c *gin.Context, db *sql.DB) {
	BasicAuth()(c)
	if c.IsAborted() {
		return
	}
	id := c.Param("ID")
	var updatedBioskop structs.Bioskop
	if err := c.ShouldBindJSON(&updatedBioskop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	if updatedBioskop.Nama == "" || updatedBioskop.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Lokasi tidak boleh kosong"})
		return
	}
	query := `UPDATE bioskop SET "Nama" = $1, "Lokasi" = $2, "Rating" = $3 WHERE "ID" = $4`
	result, err := db.Exec(query, updatedBioskop.Nama, updatedBioskop.Lokasi, updatedBioskop.Rating, id)
	if err != nil {
		log.Println("Error updating bioskop:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui bioskop"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui bioskop"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bioskop tidak ditemukan"})
		return
	}
	updatedBioskop.ID, _ = strconv.Atoi(id) // Assuming ID is int
	c.JSON(http.StatusOK, updatedBioskop)
}

// Endpoint DELETE by ID
func DeleteBioskop(c *gin.Context, db *sql.DB) {
	id := c.Param("ID")
	query := `DELETE FROM bioskop WHERE "ID" = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		log.Println("Error deleting bioskop:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus bioskop"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus bioskop"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bioskop tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil dihapus"})
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