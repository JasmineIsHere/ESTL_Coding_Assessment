package main

import (
	"awesomeProject/controllers/employees"
	"awesomeProject/daos"
	_ "awesomeProject/utils/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Service is Healthy!",
		})
	})

	employeesDAO := daos.NewEmployeesDAO()

	employees.NewHandler(employeesDAO).RouteGroup(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
