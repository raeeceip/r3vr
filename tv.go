package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/net/html"
)

type TV struct {
	*tview.Box
	content string
}

func NewTV() *TV {
	return &TV{
		Box: tview.NewBox().SetBorder(true).SetTitle("r3vr TV"),
	}
}

func (t *TV) SetContent(content string) {
	t.content = content
}

func (t *TV) Draw(screen tcell.Screen) {
	t.Box.DrawForSubclass(screen, t)
	x, y, width, height := t.Box.GetInnerRect()

	lines := strings.Split(t.content, "\n")
	for i := 0; i < height; i++ {
		if i < len(lines) {
			line := lines[i]
			tview.Print(screen, line, x, y+i, width, tview.AlignLeft, tcell.ColorWhite)
		}
	}

	// Draw TV frame
	for i := x; i < x+width; i++ {
		screen.SetContent(i, y-1, tcell.RuneHLine, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
		screen.SetContent(i, y+height, tcell.RuneHLine, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	}
	for i := y; i < y+height; i++ {
		screen.SetContent(x-1, i, tcell.RuneVLine, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
		screen.SetContent(x+width, i, tcell.RuneVLine, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	}

	// Draw TV corners
	screen.SetContent(x-1, y-1, tcell.RuneULCorner, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	screen.SetContent(x+width, y-1, tcell.RuneURCorner, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	screen.SetContent(x-1, y+height, tcell.RuneLLCorner, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	screen.SetContent(x+width, y+height, tcell.RuneLRCorner, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))

	// Draw TV antenna
	antennaX := x + width/2
	screen.SetContent(antennaX, y-2, '/', nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	screen.SetContent(antennaX+1, y-2, '\\', nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	screen.SetContent(antennaX, y-3, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	screen.SetContent(antennaX+1, y-3, '|', nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
}

func fetchPage(url string) (string, error) {
	// Create a client that doesn't follow redirects automatically
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Make the HTTP request
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching the page: %v", err)
	}
	defer resp.Body.Close()

	// Check for redirects
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location := resp.Header.Get("Location")
		return fmt.Sprintf("Redirect detected. New location: %s", location), nil
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the response: %v", err)
	}

	// Check for JavaScript redirects
	if strings.Contains(string(body), "document.location") {
		return fmt.Sprintf("JavaScript redirect detected. Content: %s", string(body)), nil
	}

	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("error parsing the page: %v", err)
	}

	// Extract text content
	var content strings.Builder
	extractText(doc, &content)

	return content.String(), nil
}

func extractText(n *html.Node, content *strings.Builder) {
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			content.WriteString(text)
			content.WriteString(" ")
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, content)
	}
	if n.Type == html.ElementNode {
		switch n.Data {
		case "p", "br", "div", "h1", "h2", "h3", "h4", "h5", "h6":
			content.WriteString("\n")
		}
	}
}
