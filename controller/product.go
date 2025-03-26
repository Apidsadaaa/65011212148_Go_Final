package controller

import (
	"go-final/dbconnect"
	"go-final/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ProductController(router *gin.Engine) {
	routes := router.Group("/product")
	{
		routes.GET("/search", searchProducts)
	}
}

func searchProducts(c *gin.Context) {
	db := dbconnect.ConnectDB()

	var products []model.Product
	var query = db.Model(&model.Product{})

	description := c.Query("description")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	// แปลงค่า minPrice และ maxPrice เป็น float64
	var minPrice, maxPrice float64
	var err error

	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price format"})
			return
		}
	}

	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_price format"})
			return
		}
	}

	// ค้นหาตาม description
	if description != "" {
		query = query.Where("description LIKE ?", "%"+description+"%")
	}

	// ค้นหาตามช่วงราคา
	if minPriceStr != "" && maxPriceStr != "" {
		query = query.Where("price BETWEEN ? AND ?", minPrice, maxPrice)
	} else if minPriceStr != "" {
		query = query.Where("price >= ?", minPrice)
	} else if maxPriceStr != "" {
		query = query.Where("price <= ?", maxPrice)
	}

	// ดึงข้อมูลสินค้า
	result := query.Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}
