package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type driver struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     uint   `json:"age"`
}

var drivers = []driver{
	{ID: "1", Name: "Harold", Surname: "Mason", Age: 52},
	{ID: "2", Name: "Leticia", Surname: "Alvarez", Age: 34},
	{ID: "3", Name: "Bojan", Surname: "Petrovic", Age: 45},
}

func main() {
	router := gin.Default()
	router.GET("/drivers", getDrivers)
	router.GET("/drivers/:id", getDriverByID)

	router.Run("localhost:8080")
}

func getDrivers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, drivers)
}

func getDriverByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := strconv.Atoi(id); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": " bad request"})
		return
	}
	for _, a := range drivers {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": " not found"})
}
