package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/pflag"
)

//go:embed icon.png
// nolint: gochecknoglobals
var icon []byte

func main() {
	windowTitle := pflag.String("title", "Enter your password", "Window title for the password prompt")

	pflag.Parse()

	app := app.New()

	icon := &fyne.StaticResource{
		StaticName:    "icon.png",
		StaticContent: icon,
	}
	app.SetIcon(icon)

	window := app.NewWindow(truncateString(*windowTitle, 25)) // TODO: arbitrary truncation
	window.SetIcon(icon)
	window.CenterOnScreen()
	window.SetPadded(true)
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(300, 50))

	var cancelled bool
	onCancel := func() {
		log.Println("Cancelled password input")

		cancelled = true
		window.Close()
	}

	window.SetCloseIntercept(onCancel)

	entry := newCancelledEntry(onCancel)
	entry.SetPlaceHolder("Enter your password...")
	entry.Password = true
	entry.Wrapping = fyne.TextTruncate
	entry.OnSubmitted = func(s string) {
		fmt.Print(s)

		window.Close()
	}
	entry.OnCancel = onCancel

	window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyEscape:
			onCancel()
		default:
			// set focus to the entry
			window.Canvas().Focus(entry)
			entry.TypedKey(key)
		}
	})

	container := container.NewVBox(entry)
	window.SetContent(container)

	window.Canvas().Focus(entry)

	window.Show()

	app.Run()

	if cancelled {
		os.Exit(1)
	}
}

type cancelledEntry struct {
	widget.Entry

	OnCancel func()
}

func newCancelledEntry(onCancel func()) *cancelledEntry {
	entry := &cancelledEntry{}

	entry.ExtendBaseWidget(entry)

	return entry
}

func (e *cancelledEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyEscape:
		if e.OnCancel != nil {
			e.OnCancel()
		} else {
			e.Entry.TypedKey(key)
		}
	default:
		e.Entry.TypedKey(key)
	}
}

func truncateString(s string, l int) string {
	if len(s) < l {
		return s
	}

	return s[:l]
}
