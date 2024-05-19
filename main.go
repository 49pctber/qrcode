package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	qrcode "github.com/skip2/go-qrcode"
	"golang.design/x/clipboard"

	_ "embed"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createQR(data string) (png []byte) {

	png, err := qrcode.Encode(data, qrcode.Medium, 500)
	checkErr(err)

	return png
}

func getClipboard() string {
	return string(clipboard.Read(clipboard.FmtText))
}

func saveQr(qrImgData []byte, win fyne.Window) {
	filedialog := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		checkErr(err)
		if uc != nil {
			defer uc.Close()
			_, err = fyne.URIWriteCloser.Write(uc, qrImgData)
			checkErr(err)
		}
	}, win)
	filedialog.SetFileName("qr.png")
	filedialog.Show()
}

func updateQr(data string, qrImgData *[]byte, qrImg *canvas.Image) {
	if data == "" {
		data = "https://bryanredd.com"
	}
	*qrImgData = createQR(data)
	qrImg.Resource = fyne.NewStaticResource("QR Code", *qrImgData)
	qrImg.Refresh()
}

//go:embed static/icon.png
var icon []byte
var resourceIconPng = &fyne.StaticResource{
	StaticName:    "icon.png",
	StaticContent: icon,
}

func main() {
	err := clipboard.Init()
	checkErr(err)

	// Create the Application
	myApp := app.New()
	myWindow := myApp.NewWindow("Instant QR")
	// icon, err := fyne.LoadResourceFromPath("icon.png")
	// checkErr(err)
	myWindow.SetIcon(resourceIconPng)

	// User Entry Box
	userEntry := widget.NewEntry()
	userEntry.SetText(getClipboard())

	// Create the QR code
	var qrImgData = createQR(userEntry.Text)
	var resource = fyne.NewStaticResource("QR Code", qrImgData)
	qrImg := canvas.NewImageFromResource(resource)
	qrImg.SetMinSize(fyne.Size{Width: 500, Height: 500}) // by default size is 0, 0

	userEntry.OnChanged = func(data string) {
		updateQr(data, &qrImgData, qrImg)
	}

	// Display a vertical box containing text, image and button
	box := container.NewVBox(
		qrImg,
		userEntry,
	)

	// Main menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Save", func() { saveQr(qrImgData, myWindow) }),
		fyne.NewMenuItem("Quit", func() { myApp.Quit() }),
	)

	mainMenu := fyne.NewMainMenu(
		fileMenu,
	)
	myWindow.SetMainMenu(mainMenu)

	// Display our content
	myWindow.SetContent(box)

	// Close the App when Escape key is pressed
	myWindow.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyEscape || keyEvent.Name == fyne.KeyQ {
			myApp.Quit()
		} else if keyEvent.Name == fyne.KeyS {
			saveQr(qrImgData, myWindow)
		} else if keyEvent.Name == fyne.KeyV {
			userEntry.SetText(getClipboard())
		} else if keyEvent.Name == fyne.KeyC {
			clipboard.Write(clipboard.FmtText, []byte(userEntry.Text))
		}
	})

	// Show window and run app
	myWindow.ShowAndRun()
}
