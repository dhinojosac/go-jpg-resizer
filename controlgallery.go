package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"strconv"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/nfnt/resize"
)

var mainwin *ui.Window
var ratio uint

func makeDataChoosersPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, true)

	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	button := ui.NewButton("Open File")
	entry := ui.NewEntry()
	entry.SetReadOnly(true)
	var filename string
	button.OnClicked(func(*ui.Button) {
		filename = ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry.SetText(filename)
	})
	grid.Append(button, 0, 0, 1, 1, false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry, 1, 0, 1, 1, true, ui.AlignFill, false, ui.AlignFill)

	labelRatio := ui.NewLabel("Ratio")
	entry2 := ui.NewEntry()
	grid.Append(labelRatio, 0, 1, 1, 1, false, ui.AlignEnd, false, ui.AlignFill)
	grid.Append(entry2, 1, 1, 1, 1, true, ui.AlignCenter, false, ui.AlignEnd)

	msggrid := ui.NewGrid()
	msggrid.SetPadded(true)
	grid.Append(msggrid,
		0, 2, 2, 1,
		false, ui.AlignCenter, false, ui.AlignStart)

	/*
		ui.MsgBoxError(mainwin,
			"This message box describes an error.",
			"More detailed information can be shown here.")
	*/

	button = ui.NewButton("Resize JPEG")
	button.OnClicked(func(*ui.Button) {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		img, err := jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()

		rat, err := strconv.Atoi(entry2.Text())
		ratio = uint(rat)
		fmt.Printf("Ratio:%v\n", ratio)
		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Resize(ratio, 0, img, resize.Lanczos3)

		out, err := os.Create("jpeg_resized.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		// write new image to file
		jpeg.Encode(out, m, nil)
	})
	msggrid.Append(button,
		1, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	return hbox
}

func setupUI() {
	mainwin = ui.NewWindow("JPG Resizer", 320, 120, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	mainwin.SetMargined(true)
	mainwin.SetChild(makeDataChoosersPage())
	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
