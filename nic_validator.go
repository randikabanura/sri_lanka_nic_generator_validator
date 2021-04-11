package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hako/durafmt"
	"net/http"
	"strconv"
	"time"
)

func validator(c *gin.Context) {
	layout := "2006-01-02"
	dns := c.Query("nic") // Data NIC string
	val18 := c.Query("val18")

	if val18 != "false" && val18 != "0" {
		val18 = "true"
	}

	if len(dns) == 0 {
		sendErrorJsonValidator(c, fmt.Errorf("nic parameter does not exist."), http.StatusBadRequest)
		return
	}

	date, doy, age, err := dateHandler(dns)

	if err != nil {
		sendErrorJsonValidator(c, err, http.StatusBadRequest)
		return
	}

	version, err := versionCheck(dns)

	if err != nil {
		sendErrorJsonValidator(c, err, http.StatusBadRequest)
		return
	}

	sex, err := sexCheck(dns)

	if err != nil {
		sendErrorJsonValidator(c, err, http.StatusBadRequest)
		return
	}

	fage, err := durafmt.ParseString(age.String())

	if err != nil {
		sendErrorJsonValidator(c, err, http.StatusBadRequest)
		return
	}

	sn, err := serialNumberHandler(dns)

	if err != nil {
		sendErrorJsonValidator(c, err, http.StatusBadRequest)
		return
	}

	cd, err := checkDigitHandler(dns)

	if err != nil {
		sendErrorJsonValidator(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"date":    date.Format(layout),
		"doy":     doy,
		"age":     fage.LimitFirstN(3).String(),
		"version": version,
		"sex":     sex,
		"sn": gin.H{
			"old": fmt.Sprintf("%03d", sn),
			"new": fmt.Sprintf("%04d", sn),
		},
		"cd":             cd,
		"validateStatus": true,
	})

}

func checkDigitHandler(dns string) (int, error) {
	if len(dns) == 10 {
		cd, errCheckDigitParse := strconv.ParseInt(string(dns[8]), 0, 64) // Check digit

		if errCheckDigitParse != nil {
			return 0, fmt.Errorf("Error occured in check digit parse.")
		}

		return int(cd), nil
	} else if len(dns) == 12 {
		cd, errCheckDigitParse := strconv.ParseInt(string(dns[len(dns)-1]), 0, 64) // Check digit

		if errCheckDigitParse != nil {
			return 0, fmt.Errorf("Error occured in check digit parse.")
		}

		return int(cd), nil
	} else {
		return 0, fmt.Errorf("Error occured on check digit handler.")
	}
}

func serialNumberHandler(dns string) (int, error) {
	if len(dns) == 10 {
		sn, errSerialNumberParse := strconv.ParseInt(dns[5:8], 10, 64) // Serial Number

		if errSerialNumberParse != nil {
			return 0, fmt.Errorf("Error occured in serial number parse.")
		}

		return int(sn), nil
	} else if len(dns) == 12 {
		sn, errSerialNumberParse := strconv.ParseInt(dns[7:11], 10, 64) // Serial Number

		if errSerialNumberParse != nil {
			return 0, fmt.Errorf("Error occured in serial number parse.")
		}

		return int(sn), nil
	} else {
		return 0, fmt.Errorf("Error occured on serial number handler.")
	}
}

func sexCheck(dns string) (string, error) {
	if len(dns) == 10 {
		sas := "Male"

		sosd, errSexParse := strconv.ParseInt(dns[2:5], 0, 64) // String of sex digits
		if errSexParse != nil {
			return "", fmt.Errorf("Error occured in day of the year parse.")
		}

		if sosd > 500 {
			sas = "Female"
		}

		return sas, nil
	} else if len(dns) == 12 {
		sas := "Male"

		sosd, errSexParse := strconv.ParseInt(dns[4:7], 0, 64) // String of sex digits
		if errSexParse != nil {
			return "", fmt.Errorf("Error occured in day of the year parse.")
		}

		if sosd > 500 {
			sas = "Female"
		}

		return sas, nil
	} else {
		return "", fmt.Errorf("Error occured on sex check.")
	}
}

func versionCheck(dns string) (string, error) {
	if len(dns) == 10 {
		return "Old", nil
	} else if len(dns) == 12 {
		return "New", nil
	} else {
		return "", fmt.Errorf("Error occured on version check.")
	}
}

func dateHandler(dns string) (time.Time, int, time.Duration, error) {
	layout := "2006-01-02"

	if len(dns) == 10 {
		year := fmt.Sprintf("19%v", dns[0:2])
		fdoy := fmt.Sprintf("%v-01-01", year)
		date, errTimeParse := time.Parse(layout, fdoy)
		if errTimeParse != nil {
			return time.Now(), 0, 0, fmt.Errorf("Error occured in year parse.")
		}

		doy, errIntParse := strconv.ParseInt(dns[2:5], 0, 64)
		if errIntParse != nil {
			return time.Now(), 0, 0, fmt.Errorf("Error occured in day of the year parse.")
		}

		if doy > 500 {
			doy -= 500
		}
		date = date.AddDate(0, 0, int(doy-1))
		age := time.Since(date)

		return date, int(doy), age, nil
	} else if len(dns) == 12 {
		year := dns[0:4]
		fdoy := fmt.Sprintf("%v-01-01", year)
		date, errTimeParse := time.Parse(layout, fdoy)
		if errTimeParse != nil {
			return time.Now(), 0, 0, fmt.Errorf("Error occured in year parse.")
		}

		doy, errIntParse := strconv.ParseInt(dns[4:7], 0, 64)
		if errIntParse != nil {
			return time.Now(), 0, 0, fmt.Errorf("Error occured in day of the year parse.")
		}

		if doy > 500 {
			doy -= 500
		}
		date = date.AddDate(0, 0, int(doy-1))
		age := time.Since(date)

		return date, int(doy), age, nil
	} else {
		return time.Now(), 0, 0, fmt.Errorf("nic parameter value is incorrect.")
	}
}

//func validate18Years(val18 string) {
//	//bval18 := true
//	//layout := "2006-01-02"
//	//
//	//if val18 == "false" || val18 == "0" {
//	//	bval18 = false
//	//}
//}

func sendErrorJsonValidator(c *gin.Context, err error, code int) {
	c.JSON(code, gin.H{
		"status":         false,
		"error":          err.Error(),
		"code":           http.StatusText(code),
		"validateStatus": false,
	})
}
