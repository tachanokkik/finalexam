package customer

import (
	"github.com/tachanokkik/finalexam/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func createCustomerHandler(c *gin.Context) {
	customer := database.Customer{}
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	conn := database.Conn()
	err := conn.CreateCustomer(&customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func getCustomerByIdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	conn := database.Conn()
	customer, err := conn.GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func getCustomerAllHandler(c *gin.Context) {
	conn := database.Conn()
	customers, err := conn.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, customers)
}

func updateCustomerByIdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	customer := database.Customer{}
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	conn := database.Conn()
	if err := conn.UpdateById(id, &customer); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func deleteCustomerByIdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	conn := database.Conn()
	if err := conn.DeleteById(id); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, gin.H{"message":"customer deleted"})
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(AuthMiddleware)

	r.POST("/customers", createCustomerHandler)
	r.GET("/customers/:id", getCustomerByIdHandler)
	r.GET("/customers", getCustomerAllHandler)
	r.PUT("/customers/:id", updateCustomerByIdHandler)
	r.DELETE("/customers/:id", deleteCustomerByIdHandler)

	return r
}
