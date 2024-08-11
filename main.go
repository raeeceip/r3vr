// main.go
package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	tv := NewTV()

	urlInput := tview.NewInputField().
		SetLabel("URL: ").
		SetFieldWidth(50)

	urlInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			url := urlInput.GetText()
			go func() {
				content, err := fetchPage(url)
				if err != nil {
					content = err.Error()
				}
				app.QueueUpdateDraw(func() {
					tv.SetContent(content)
				})
			}()
		}
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(urlInput, 1, 1, true).
		AddItem(tv, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		log.Fatal(err)
	}
}
