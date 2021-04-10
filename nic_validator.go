package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func validator(c *gin.Context) {
	//layout := "2006-01-02"
	dns := c.Query("nic")
	if len(dns) == 0 {
		sendErrorJsonValidator(c, fmt.Errorf("nic parameter does not exist."))
	}

}

func sendErrorJsonValidator(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status": false,
		"error":  err.Error(),
	})
}
