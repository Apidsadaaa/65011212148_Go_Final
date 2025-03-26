package controller

import (
	"fmt"
	"go-final/dbconnect"
	"go-final/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CustomerController(router *gin.Engine) {
	routes := router.Group("/customer")
	{
		routes.GET("/", getAllCustomer)
		routes.POST("/login", loginCustomer)
		routes.PUT("/change", changePassword)
	}
}

func getAllCustomer(c *gin.Context) {
	db := dbconnect.ConnectDB() // เชื่อมต่อฐานข้อมูล

	var customers []model.Customer
	result := db.Find(&customers)

	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"customers": customers})
}

func loginCustomer(c *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	db := dbconnect.ConnectDB()

	var customer model.Customer
	// ค้นหาลูกค้าจาก email
	result := db.Where("email = ?", request.Email).First(&customer)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	// ตรวจสอบรหัสผ่าน
	if customer.Password != request.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}
	// ส่งข้อมูลลูกค้าที่ผ่านการตรวจสอบแล้ว
	c.JSON(http.StatusOK, gin.H{"customer": customer})
}
func changePassword(c *gin.Context) {
	var request struct {
		Email           string `json:"email"`
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	// รับค่าจาก request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if request.NewPassword != request.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password and confirm password do not match"})
		return
	}

	db := dbconnect.ConnectDB()

	var customer model.Customer
	result := db.Where("email = ?", request.Email).First(&customer)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// ตรวจสอบรหัสผ่านเก่า
	if customer.Password != request.OldPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
		return
	}

	// อัปเดตรหัสผ่านใหม่
	customer.Password = request.NewPassword
	db.Save(&customer)

	// ส่งข้อมูลที่สำเร็จกลับไป
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
