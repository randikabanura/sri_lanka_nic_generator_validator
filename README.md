# Sri Lanka NIC Generator and Validator

This is a project to implement a Generator and Validator for Sri Lanka National Identity Card (NIC).
This is currently an API only project build using a state-of-the-art programming language of [Go](https://golang.org/).
This project also utilize Few Go packages including:
  * [Gin Web Framework](https://github.com/gin-gonic/) - For API implementation
  * [durafmt](https://github.com/hako/durafmt) - For getting duration between Times
  * [go-randomdata](https://github.com/Pallinder/go-randomdata) - For getting random Int, Boolean, Times
  * [barcode](https://github.com/boombuler/barcode) - For generating the PDF417 barcode that in the NIC version introduced in 2016

This is currently ongoing project and hope to improve the functionality lot more.

## Usage

### NIC generator

For generation of and NIC call ```/v1/generator``` endpoint. It has a few query params it could take
to generate NIC number according to the parameters. These params includes the following:
  * sex - f, female, m , male
  * date - date of which NIC holder is born (Ex: 1995-05-17)

#### Provinces

API is returning a province which includes a number and a name. This province name or number does not
relate to NIC number itself. This is given because Old NIC type had a number indicating the province in the 
physical NIC. List of the provinces, and the numbers as follows:

  1. Western Province
  2. Central Province
  3. Southern Province
  4. Northern Province
  5. Eastern Province
  6. North Western Province
  7. North Central Province
  8. Uva Province
  9. Sabaragahmuwa Province

#### Barcodes

New NIC card has a barcode which contain the same data that written in the card itself.
This barcode is type of PDF417 and can be scanned and get the data using simple barcode scanning app.
I have implemented the generation of barcode as a base64 image that can be transferred in an API call.
Generator end point response contain the barcode object for the content nad image of the generated data.
Sample barcode image and data contain in that given below. 

```
00
196165906402
08/06/1961
Female
1972-05-29
00BFT-710
Ava Martinez
92 Madison Pkwy,
Northleach, NC, 35229
Ransom Canyon
485406C0548FDD8FDDF300F312EE947D#
```
![img.png](sample_nic_barcode.png)

#### Examples

```
// http://localhost:3000/v1/generator?sex=male
// Usage of the sex param

{
  "cd": 4,
  "date": "1948-12-23",
  "doy": 358,
  "nnic": "19483580944",
  "onic": "48358944V",
  "province": {
    "name": "Western Province",
    "number": 1
  },
  "sex": "Male",
  "sn": {
    "new": "0094",
    "old": "094"
  },
  "status": true,
  "barcode": {
    "content": "00\n199427605649\n03/10/1994\nMale\n2005-03-20.........",
    "image": "data:image/png;base64,iVBORw0...."
  }
}

// http://localhost:3000/v1/generator?date=1995-05-17
// Usage of the date param

{
  "cd": 4,
  "date": "1995-05-17",
  "doy": 137,
  "nnic": "199563701074",
  "onic": "956371074V",
  "province": {
    "name": "Southern Province",
    "number": 3
  },
  "sex": "Female",
  "sn": {
    "new": "0107",
    "old": "107"
  },
  "status": true,
  "barcode": {
    "content": "00\n199427605649\n03/10/1994\nMale\n2005-03-20.........",
    "image": "data:image/png;base64,iVBORw0...."
  }
}

// http://localhost:3000/v1/generator?date=1996-06-03&sex=m
// Usage of the both sex and date param

{
  "cd": 6,
  "date": "1996-06-03",
  "doy": 155,
  "nnic": "199615506166",
  "onic": "961556166V",
  "province": {
    "name": "Western Province",
    "number": 1
  },
  "sex": "Male",
  "sn": {
    "new": "0616",
    "old": "616"
  },
  "status": true,
  "barcode": {
    "content": "00\n199427605649\n03/10/1994\nMale\n2005-03-20.........",
    "image": "data:image/png;base64,iVBORw0...."
  }
}

// http://localhost:3000/v1/generator?date=2005-05-17
// Usage of if birthday is in or after year 2000
// Note that if the birthday is in or after year 2000,
// response will not have onic (Old NIC) and Old serial number

{
  "cd": 4,
  "date": "2005-05-17",
  "doy": 137,
  "nnic": "200563797534",
  "onic": "",
  "sex": "Female",
  "sn": {
    "new": "9753",
    "old": ""
  },
  "status": true,
  "barcode": {
    "content": "00\n199427605649\n03/10/1994\nMale\n2005-03-20.........",
    "image": "data:image/png;base64,iVBORw0...."
  }
}
```

### NIC validator
For validation of the NIC call ```/v1/validator``` endpoint. It must have a ```nic``` parameter.
Otherwise it will return an error saying that ```nic``` parameter is empty. NIC that send in the parameter ```nic``` could be
both new of old version of the NIC in sri lanka

#### Examples

```
// http://localhost:3000/v1/validator?nic=956380995V
// Using old version of the NIC

{
  "age": "25 years 47 weeks 6 days",
  "cd": 5,
  "date": "1995-05-18",
  "doy": 138,
  "sex": "Female",
  "status": true,
  "sn": {
    "new": "0099",
    "old": "099"
  },
  "validateStatus": true,
  "version": "Old"
}

// http://localhost:3000/v1/validator?nic=199615500343
// Using new version of the NIC

{
  "age": "24 years 45 weeks 3 days",
  "cd": 3,
  "date": "1996-06-03",
  "doy": 155,
  "sex": "Male",
  "sn": {
    "new": "0034",
    "old": "034"
  },
  "status": true,
  "validateStatus": true,
  "version": "New"
}

// http://localhost:3000/v1/validator?nic=9563809s95V
// Response when NIC is not valid

{
  "code": "Bad Request",
  "error": "nic parameter value is incorrect.",
  "status": false,
  "validateStatus": false
}

```

## Do you like it? Star it!
If you use this component just star it. A developer is more motivated to improve a project when there is some interest.

## Contributing
Bug reports and pull requests are welcome on GitHub at https://github.com/randikabanura/sri-lanka-nic-generator-validator.

## License
The software is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).

## Developer
Name: [Banura Randika Perera](https://github.com/randikabanura) <br/>
Linkedin: [randika-banura](https://www.linkedin.com/in/randika-banura/) <br/>
Email: [randika.banura@gamil.com](mailto:randika.banura@gamil.com) <br/>
Bsc (Hons) Information Technology specialized in Software Engineering (SLIIT)