package main

import (
	"encoding/base64"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/boombuler/barcode/pdf417"
	"github.com/gin-gonic/gin"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	"os"
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
	doy := int(math.Ceil(date.Sub(fdoy).Hours()/24) + 1)       // Day of the year

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
	barcode := ""
	barcodeContent := ""
	if len(nnic) != 12 {
		nnic = ""
		nsn = ""
	} else {
		barcodeContent, barcode, err = generateBarcodeForNNIC(nnic, date, sas)
		if err != nil {
			sendErrorJsonGenerator(c, err, http.StatusBadRequest)
			return
		}
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
		"barcode": gin.H{
			"content": barcodeContent,
			"image":   barcode,
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
func generateONIC(year int, doy int, sn int, cd int) string {
	if sn > 999 || year > 2000 {
		return ""
	}

	sy := year % 100
	ssy := fmt.Sprintf("%v", sy)

	if sy < 10 {
		ssy = fmt.Sprintf("0%v", sy)
	}

	return fmt.Sprintf("%v%03d%03d%d%v", ssy, doy, sn, cd, "V")
}

// Generate new nic version according to year, day of the year, serial number and check digit
func generateNNIC(year int, doy int, sn int, cd int) string {
	return fmt.Sprintf("%d%03d%04d%d", year, doy, sn, cd)
}

// Generate the pdf417 barcode for the new NIC number
func generateBarcodeForNNIC(nnic string, date time.Time, sas string) (string, string, error) {
	layoutOne := "2006-01-02"
	layoutTwo := "02/01/2006"
	fullName := ""
	if sas == "Male" {
		fullName = randomdata.FullName(randomdata.Male)
	} else {
		fullName = randomdata.FullName(randomdata.Female)
	}

	db := date.Format(layoutOne)
	db16 := date.AddDate(16, 0, 0).Format(layoutOne) // Date 16 years after today
	createdDate, err := time.Parse("Monday 2 Jan 2006", randomdata.FullDateInRange(db, db16))
	if err != nil {
		return "", "", err
	}

	barcodeStr := "00\n" +
		nnic + "\n" +
		date.Format(layoutTwo) + "\n" +
		sas + "\n" +
		createdDate.Format(layoutOne) + "\n" +
		"00BFT-710\n" +
		fullName + "\n" +
		randomdata.Address() + "\n" +
		randomdata.City() + "\n" +
		"485406C0548FDD8FDDF300F312EE947D#"

	currentTime := time.Now().Unix()

	pdf417Code, err := pdf417.Encode(barcodeStr, 0)
	if err != nil {
		return "", "", err
	}

	// create the output file
	file, err := os.Create(fmt.Sprintf("barcode-%v.png", currentTime))
	if err != nil {
		return "", "", err
	}

	defer file.Close()

	err = png.Encode(file, pdf417Code)
	if err != nil {
		return "", "", err
	}

	bytes, err := ioutil.ReadFile(fmt.Sprintf("./barcode-%v.png", currentTime))
	if err != nil {
		return "", "", err
	}

	err = os.Remove(fmt.Sprintf("barcode-%v.png", currentTime))
	if err != nil {
		return "", "", err
	}

	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	base64Encoding := getMimeType(mimeType)

	// Append the base64 encoded output
	base64Encoding += toBase64(bytes)

	// Print the full base64 representation of the image
	return pdf417Code.Content(), base64Encoding, nil
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func getMimeType(mimetype string) string {
	base64Encoding := ""

	switch mimetype {
	case "image/jpeg":
		base64Encoding = "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding = "data:image/png;base64,"
	}

	return base64Encoding
}
