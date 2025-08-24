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

// GET /api/books
func GetBooks(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, title, category_id, description, image_url, release_year, price, total_page, thickness, created_at, created_by, modified_at, modified_by
		FROM buku`)
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

		if err := rows.Scan(&id, &b.Title, &b.CategoryId, &b.Description, &b.ImageUrl, &b.ReleaseYear, &b.Price, &b.TotalPage, &b.Thickness, &b.CreatedAt, &b.CreatedBy, &modifiedAt, &modifiedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data buku", "message": err.Error()})
			return
		}

		if modifiedAt.Valid {
			b.ModifiedAt = &modifiedAt.Time
		}
		if modifiedBy.Valid {
			b.ModifiedBy = &modifiedBy.String
		}

		books = append(books, b)
	}

	c.JSON(http.StatusOK, books)
}

// POST /api/books
func CreateBook(c *gin.Context) {
	var input struct {
		Title       string `json:"title"`
		CategoryId  int    `json:"categoryId"`
		Description string `json:"description"`
		ImageUrl    string `json:"imageUrl"`
		ReleaseYear int    `json:"releaseYear"`
		Price       int    `json:"price"`
		TotalPage   int    `json:"totalPage"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "message": err.Error()})
		return
	}

	if input.ReleaseYear < 1980 || input.ReleaseYear > 2024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Release year harus antara 1980-2024"})
		return
	}

	thickness := "tipis"
	if input.TotalPage > 100 {
		thickness = "tebal"
	}

	var createdBy string
	if user, exists := c.Get("username"); exists {
		createdBy = user.(string)
	} else {
		createdBy = "unknown"
	}

	createdAt := time.Now()

	sqlStatement := `
		INSERT INTO buku (title, category_id, description, image_url, release_year, price, total_page, thickness, created_at, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id
	`
	var id int
	err := database.DB.QueryRow(sqlStatement, input.Title, input.CategoryId, input.Description, input.ImageUrl, input.ReleaseYear, input.Price, input.TotalPage, thickness, createdAt, createdBy).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan buku", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          id,
		"title":       input.Title,
		"categoryId":  input.CategoryId,
		"description": input.Description,
		"imageUrl":    input.ImageUrl,
		"releaseYear": input.ReleaseYear,
		"price":       input.Price,
		"totalPage":   input.TotalPage,
		"thickness":   thickness,
		"createdAt":   createdAt,
		"createdBy":   createdBy,
	})
}

// GET /api/books/:id
func GetBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var b models.Book
	var modifiedAt sql.NullTime
	var modifiedBy sql.NullString

	err := database.DB.QueryRow(`
		SELECT title, category_id, description, image_url, release_year, price, total_page, thickness, created_at, created_by, modified_at, modified_by
		FROM buku WHERE id=$1`, id).Scan(&b.Title, &b.CategoryId, &b.Description, &b.ImageUrl, &b.ReleaseYear, &b.Price, &b.TotalPage, &b.Thickness, &b.CreatedAt, &b.CreatedBy, &modifiedAt, &modifiedBy)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	if modifiedAt.Valid {
		b.ModifiedAt = &modifiedAt.Time
	}
	if modifiedBy.Valid {
		b.ModifiedBy = &modifiedBy.String
	}

	c.JSON(http.StatusOK, b)
}

// DELETE /api/books/:id
func DeleteBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	res, err := database.DB.Exec(`DELETE FROM buku WHERE id=$1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus buku", "message": err.Error()})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dihapus"})
}

// PUT /api/books/:id
func UpdateBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Title       string `json:"title"`
		CategoryId  int    `json:"categoryId"`
		Description string `json:"description"`
		ImageUrl    string `json:"imageUrl"`
		ReleaseYear int    `json:"releaseYear"`
		Price       int    `json:"price"`
		TotalPage   int    `json:"totalPage"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "message": err.Error()})
		return
	}

	if input.ReleaseYear < 1980 || input.ReleaseYear > 2024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Release year harus antara 1980-2024"})
		return
	}

	thickness := "tipis"
	if input.TotalPage > 100 {
		thickness = "tebal"
	}

	var modifiedBy string
	if user, exists := c.Get("username"); exists {
		modifiedBy = user.(string)
	} else {
		modifiedBy = "unknown"
	}
	modifiedAt := time.Now()

	res, err := database.DB.Exec(`
		UPDATE buku
		SET title=$1, category_id=$2, description=$3, image_url=$4, release_year=$5, price=$6, total_page=$7, thickness=$8, modified_at=$9, modified_by=$10
		WHERE id=$11`,
		input.Title, input.CategoryId, input.Description, input.ImageUrl, input.ReleaseYear, input.Price, input.TotalPage, thickness, modifiedAt, modifiedBy, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update buku", "message": err.Error()})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Buku berhasil diperbarui",
		"id":          id,
		"title":       input.Title,
		"categoryId":  input.CategoryId,
		"description": input.Description,
		"imageUrl":    input.ImageUrl,
		"releaseYear": input.ReleaseYear,
		"price":       input.Price,
		"totalPage":   input.TotalPage,
		"thickness":   thickness,
		"modifiedAt":  modifiedAt,
		"modifiedBy":  modifiedBy,
	})
}
