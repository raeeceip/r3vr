// main.go
package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	browser := NewBrowser()

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
					browser.SetContent(content)
				})
			}()
		}
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(urlInput, 1, 1, true).
		AddItem(browser, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		fmt.Println("Error running application:", err)
	}
}
