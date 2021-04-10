package main

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()

	r.GET("/", ping)
	r.GET("/ping", ping)
	r.GET("/generator", generator)
	r.Run(":3000")
}

func ping(c *gin.Context)  {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func generator(c *gin.Context)  {
	layout := "2006-01-02"
	qs := c.Query("date")
	date, err := queryHandler(qs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error": err,
		})
		return
	}

	fdoy := time.Date(date.Year(), 1,1,0,0,0,0,time.UTC)
	doe := math.Ceil(date.Sub(fdoy).Hours() / 24) + 1

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"date": date.Format(layout),
		"doe": doe,
	})
}

func queryHandler(ds string) (time.Time, error) {
	layout := "2006-01-02"
	date := time.Now()
	var err error = nil

	if len(ds) > 0 {
		date, err = time.Parse(layout, ds)
	} else {
		db18 := time.Now().AddDate(-18, 0, 0).Format(layout)
		db118 := time.Now().AddDate(-118, 0, 0).Format(layout)
		date, err = time.Parse("Monday 2 Jan 2006", randomdata.FullDateInRange(db118, db18))
	}

	return date, err
}