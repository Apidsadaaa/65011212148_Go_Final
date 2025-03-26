package controller

import (
	"go-final/dbconnect"
	"go-final/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CartController(router *gin.Engine) {
	routes := router.Group("/cart")
	{
		routes.POST("/add", addToCart)
		routes.GET("/view", viewCart)
	}
}

func addToCart(c *gin.Context) {
	var request struct {
		CustomerID int    `json:"customer_id"`
		CartName   string `json:"cart_name"`
		ProductID  int    `json:"product_id"`
		Quantity   int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	db := dbconnect.ConnectDB()

	var cart model.Cart
	result := db.Where("customer_id = ? AND cart_name = ?", request.CustomerID, request.CartName).First(&cart)
	if result.Error != nil {

		cart = model.Cart{
			CustomerID: request.CustomerID,
			CartName:   request.CartName,
		}
		db.Create(&cart)
	}

	var cartItem model.CartItem
	result = db.Where("cart_id = ? AND product_id = ?", cart.CartID, request.ProductID).First(&cartItem)
	if result.Error != nil {
		// หากไม่พบสินค้าในรถเข็น ให้เพิ่มสินค้าลงในรถเข็น
		cartItem = model.CartItem{
			CartID:    cart.CartID,
			ProductID: request.ProductID,
			Quantity:  request.Quantity,
		}
		db.Create(&cartItem)
	} else {
		// หากสินค้าพบในรถเข็นแล้ว ให้เพิ่มจำนวนสินค้า
		cartItem.Quantity += request.Quantity
		db.Save(&cartItem)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
}

func viewCart(c *gin.Context) {
	customerID := c.DefaultQuery("customer_id", "")

	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	db := dbconnect.ConnectDB()

	var carts []model.Cart
	result := db.Where("customer_id = ?", customerID).Find(&carts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carts"})
		return
	}

	// เตรียมข้อมูลผลลัพธ์ที่จะส่งกลับ
	var cartDetails []gin.H

	// ลูปผ่านรถเข็นแต่ละคัน
	for _, cart := range carts {
		var cartItems []model.CartItem
		db.Where("cart_id = ?", cart.CartID).Find(&cartItems)

		var itemsDetail []gin.H
		// ลูปผ่านสินค้าทุกชิ้นในรถเข็น
		for _, cartItem := range cartItems {
			var product model.Product
			db.Where("product_id = ?", cartItem.ProductID).First(&product)

			// เพิ่มข้อมูลของสินค้าลงในรายการ
			itemsDetail = append(itemsDetail, gin.H{
				"product_name": product.ProductName,
				"quantity":     cartItem.Quantity,
				"price":        product.Price,
			})
		}

		// เพิ่มข้อมูลของรถเข็นแต่ละคัน
		cartDetails = append(cartDetails, gin.H{
			"cart_name": cart.CartName,
			"items":     itemsDetail,
		})
	}

	c.JSON(http.StatusOK, gin.H{"carts": cartDetails})
}
