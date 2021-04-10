package main

import (
	"fmt"
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

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func generator(c *gin.Context) {
	layout := "2006-01-02"
	qs := c.Query("date")
	date, err := queryHandler(qs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	fdoy := time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	doy := math.Ceil(date.Sub(fdoy).Hours()/24) + 1
	sn := randomdata.Number(0, 1000)
	cd := randomdata.Number(0, 10)
	sex := randomdata.Boolean()
	sas := "Male"

	if sex == false {
		doy += 500
		sas = "Female"
	}

	onic := generateONIC(date.Year(), doy, sn, cd)
	nnic := generateNNIC(date.Year(), doy, sn, cd)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"date":   date.Format(layout),
		"doy":    doy,
		"sn":     sn,
		"cd":     cd,
		"sex":    sas,
		"onic":   onic,
		"nnic":   nnic,
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

func generateONIC(year int, doy float64, sn int, cd int) string {
	sy := year % 100
	ssy := fmt.Sprintf("%v", sy)

	if sy < 10 {
		ssy = fmt.Sprintf("0%v", sy)
	}

	return fmt.Sprintf("%v%.0f%d%d%v", ssy, doy, sn, cd, "V")
}

func generateNNIC(year int, doy float64, sn int, cd int) string {
	return fmt.Sprintf("%d%.0f0%d%d", year, doy, sn, cd)
}
