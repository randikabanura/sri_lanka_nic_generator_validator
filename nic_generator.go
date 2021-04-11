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
	provinces := []string{"Western", "Central", "Southern", "Northern", "Eastern", "North Western", "North Central", "Uva", "Sabaragahmuwa"}

	date, err := dateQueryHandler(dqs) // Date query string

	if err != nil {
		sendErrorJsonGenerator(c, err, http.StatusBadRequest)
		return
	}

	sqs := c.Query("sex") // Sex query string
	sex, sas, err := sexQueryHandler(sqs)

	if err != nil {
		sendErrorJsonGenerator(c, err, http.StatusBadRequest)
		return
	}

	fdoy := time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.UTC) // First day of the year
	doy := math.Ceil(date.Sub(fdoy).Hours()/24) + 1            // Day of the year

	sn := randomdata.Number(0, 1000) // Serial number
	if date.Year() >= 2000 {
		sn = randomdata.Number(0, 10000)
	}

	cd := randomdata.Number(0, 10) // Check digit

	sdoy := doy // Day of the year according to sex
	if sex == false {
		sdoy += 500
	}

	onic := generateONIC(date.Year(), sdoy, sn, cd)
	osn := fmt.Sprintf("%03d", sn) // Old serial number
	if len(onic) != 10 {
		onic = ""
		osn = ""
	}

	nnic := generateNNIC(date.Year(), sdoy, sn, cd)
	nsn := fmt.Sprintf("%04d", sn) // New serial number
	if len(nnic) != 12 {
		nnic = ""
		nsn = ""
	}

	pn := randomdata.Number(0, 9)                   // Province number. This is not associated with the NIC number
	ps := fmt.Sprintf("%v Province", provinces[pn]) // String of province

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"date":   date.Format(layout),
		"doy":    doy,
		"sn": gin.H{
			"old": osn, // Old serial number
			"new": nsn, // New serial number
		},
		"cd":   cd,
		"sex":  sas,
		"onic": onic, // Old nic version
		"nnic": nnic, // New nic version
		"province": gin.H{ // Province is not associated with the NIC number
			"number": pn + 1,
			"name":   ps,
		},
	})
}

// Handle error response if any error occurred for generator
func sendErrorJsonGenerator(c *gin.Context, err error, code int) {
	c.JSON(code, gin.H{
		"status": false,
		"error":  err.Error(),
		"code":   http.StatusText(code),
	})
}

// Handles a date query param. If not available it auto generate random date
func dateQueryHandler(dqs string) (time.Time, error) {
	layout := "2006-01-02"
	date := time.Now()
	var err error = nil

	if len(dqs) > 0 {
		date, err = time.Parse(layout, dqs)
	} else {
		db18 := time.Now().AddDate(-18, 0, 0).Format(layout)   // Date 18 years before today
		db118 := time.Now().AddDate(-118, 0, 0).Format(layout) // Date 118 years before today
		date, err = time.Parse("Monday 2 Jan 2006", randomdata.FullDateInRange(db118, db18))
	}

	return date, err
}

// Handles the sex query param and return a boolean, string and error
func sexQueryHandler(sqs string) (bool, string, error) {
	sqs = strings.ToLower(sqs)

	switch sqs {
	case "m", "male":
		{
			return true, "Male", nil
		}
	case "f", "female":
		{
			return false, "Female", nil
		}
	case "":
		{
			rs := randomdata.Boolean() // Random sex boolean
			rss := "Male"              // initialize the the sex string

			// Change the sex string if rs is false
			if rs == false {
				rss = "Female"
			}
			return rs, rss, nil
		}
	default:
		return false, "", fmt.Errorf("Sex parameter can not be parsed.")
	}
}

// Generate old nic version according to year, day of the year, serial number and check digit
func generateONIC(year int, doy float64, sn int, cd int) string {
	if sn > 999 {
		return ""
	}

	sy := year % 100
	ssy := fmt.Sprintf("%v", sy)

	if sy < 10 {
		ssy = fmt.Sprintf("0%v", sy)
	}

	return fmt.Sprintf("%v%.0f%03d%d%v", ssy, doy, sn, cd, "V")
}

// Generate new nic version according to year, day of the year, serial number and check digit
func generateNNIC(year int, doy float64, sn int, cd int) string {
	return fmt.Sprintf("%d%.0f%04d%d", year, doy, sn, cd)
}
