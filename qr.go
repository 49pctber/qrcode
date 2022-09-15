package main

import (
	qrcode "github.com/skip2/go-qrcode"
)

func createQR(data string) (png []byte) {

	png, err := qrcode.Encode(data, qrcode.Medium, 500)
	checkErr(err)

	// f, err := os.Create("qr.png")
	// checkErr(err)
	// defer f.Close()

	// _, err = f.Write(png)
	// checkErr(err)

	return png
}
