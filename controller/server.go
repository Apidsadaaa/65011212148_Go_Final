package controller

import "github.com/gin-gonic/gin"

func StartServer() {

	router := gin.Default()
	CustomerController(router)
	router.Run()
}
