package main

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strings"
	"time"
)

func generator(c *gin.Context) {
	layout := "2006-01-02"
	dqs := c.Query("date")
	date, err := dateQueryHandler(dqs)

	if err != nil {
		sendErrorJsonGenerator(c, err, http.StatusBadRequest)
		return
	}

	sqs := c.Query("sex")
	sex, sas, err := sexQueryHandler(sqs)

	if err != nil {
		sendErrorJsonGenerator(c, err, http.StatusBadRequest)
		return
	}

	fdoy := time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	doy := math.Ceil(date.Sub(fdoy).Hours()/24) + 1
	sn := randomdata.Number(0, 1000)
	cd := randomdata.Number(0, 10)

	if sex == false {
		doy += 500
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

func sendErrorJsonGenerator(c *gin.Context, err error, code int) {
	c.JSON(code, gin.H{
		"status": false,
		"error":  err.Error(),
		"code":   http.StatusText(code),
	})
}

func dateQueryHandler(dqs string) (time.Time, error) {
	layout := "2006-01-02"
	date := time.Now()
	var err error = nil

	if len(dqs) > 0 {
		date, err = time.Parse(layout, dqs)
	} else {
		db18 := time.Now().AddDate(-18, 0, 0).Format(layout)
		db118 := time.Now().AddDate(-118, 0, 0).Format(layout)
		date, err = time.Parse("Monday 2 Jan 2006", randomdata.FullDateInRange(db118, db18))
	}

	return date, err
}

func sexQueryHandler(sqs string) (bool, string, error) {
	sqs = strings.ToLower(sqs)

	switch sqs {
	case "m":
		{
			return true, "Male", nil
		}
	case "male":
		{
			return true, "Male", nil
		}
	case "f":
		{
			return false, "Female", nil
		}
	case "female":
		{
			return false, "Female", nil
		}
	case "":
		{
			rs := randomdata.Boolean()
			rss := "Male"

			if rs == false {
				rss = "Female"
			}
			return rs, rss, nil
		}
	default:
		return false, "", fmt.Errorf("Sex parameter can not be parsed.")
	}
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
