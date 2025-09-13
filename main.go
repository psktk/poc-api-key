package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const allowedIP = "127.0.0.1"
const apiKeyHeader = "X-API-Key"
const validAPIKey = "supersecretkey"

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	// Initialize SQLite DB
	db, err := sql.Open("sqlite3", "./products.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS product (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		price INTEGER
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert sample records if table is empty
	row := db.QueryRow("SELECT COUNT(*) FROM product")
	var count int
	row.Scan(&count)
	if count == 0 {
		_, err = db.Exec("INSERT INTO product (name, price) VALUES (?, ?), (?, ?)", "Apple", 100, "Banana", 50)
		if err != nil {
			log.Fatal(err)
		}
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "::1" {
			ip = "127.0.0.1"
		}
		if ip == allowedIP {
			c.Request.Header.Set("X-API-Key", validAPIKey)
		}
		c.Next()
	})
	r.Use(func(c *gin.Context) {
		key := c.GetHeader(apiKeyHeader)
		if key != validAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}
		c.Next()
	})

	r.POST("/products", func(c *gin.Context) {
		var p Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := db.Exec("INSERT INTO product (name, price) VALUES (?, ?)", p.Name, p.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		id, _ := res.LastInsertId()
		p.ID = int(id)
		c.JSON(http.StatusCreated, p)
	})

	r.GET("/products", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, price FROM product")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			products = append(products, p)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		row := db.QueryRow("SELECT id, name, price FROM product WHERE id = ?", id)
		var p Product
		if err := row.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	r.PUT("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var p Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := db.Exec("UPDATE product SET name = ?, price = ? WHERE id = ?", p.Name, p.Price, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		p.ID = atoi(id)
		c.JSON(http.StatusOK, p)
	})

	r.DELETE("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM product WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	r.Run()
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
