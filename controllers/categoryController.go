package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"quiz3-go/database"
	"quiz3-go/models"

	"github.com/gin-gonic/gin"
)

// GET /api/categories
func GetCategories(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, created_at, created_by, modified_at, modified_by FROM kategori")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil kategori"})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		var id int
		if err := rows.Scan(&id, &cat.Name, &cat.CreatedAt, &cat.CreatedBy, &cat.ModifiedAt, &cat.ModifiedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data kategori", "message": err.Error()})
			return
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, categories)
}

// POST /api/categories
func CreateCategory(c *gin.Context) {
	var input models.Category
	if err := c.ShouldBindJSON(&input); err != nil || input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori harus diisi"})
		return
	}

	input.CreatedAt = time.Now()
	if user, exists := c.Get("username"); exists {
		input.CreatedBy = user.(string)
	} else {
		input.CreatedBy = "unknown"
	}

	sqlStatement := `INSERT INTO kategori (name, created_at, created_by) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := database.DB.QueryRow(sqlStatement, input.Name, input.CreatedAt, input.CreatedBy).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        id,
		"title":     input.Name,
		"createdAt": input.CreatedAt,
		"createdBy": input.CreatedBy,
	})
}

// GET /api/categories/:id
func GetCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var cat models.Category
	sqlStatement := `SELECT name, created_at, created_by, modified_at, modified_by FROM kategori WHERE id=$1`
	err := database.DB.QueryRow(sqlStatement, id).Scan(&cat.Name, &cat.CreatedAt, &cat.CreatedBy, &cat.ModifiedAt, &cat.ModifiedBy)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

// DELETE /api/categories/:id
func DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	sqlStatement := `DELETE FROM kategori WHERE id=$1`
	res, err := database.DB.Exec(sqlStatement, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}

func GetBooksByCategory(c *gin.Context) {
	categoryId, _ := strconv.Atoi(c.Param("id"))

	// Cek apakah kategori ada
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM kategori WHERE id=$1)", categoryId).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengecek kategori", "message": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	// Ambil semua buku berdasarkan category_id
	rows, err := database.DB.Query(`
		SELECT id, title, description, image_url, release_year, price, total_page, thickness, created_at, created_by, modified_at, modified_by
		FROM buku
		WHERE category_id=$1`, categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil buku", "message": err.Error()})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		var modifiedAt sql.NullTime
		var modifiedBy sql.NullString
		var id int

		if err := rows.Scan(&id, &b.Title, &b.Description, &b.ImageUrl, &b.ReleaseYear, &b.Price, &b.TotalPage, &b.Thickness, &b.CreatedAt, &b.CreatedBy, &modifiedAt, &modifiedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data buku", "message": err.Error()})
			return
		}

		if modifiedAt.Valid {
			b.ModifiedAt = &modifiedAt.Time
		}
		if modifiedBy.Valid {
			b.ModifiedBy = &modifiedBy.String
		}

		b.CategoryId = categoryId
		books = append(books, b)
	}

	c.JSON(http.StatusOK, gin.H{
		"categoryId": categoryId,
		"books":      books,
	})
}
