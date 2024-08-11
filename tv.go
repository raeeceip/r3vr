package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/net/html"
)

type Element struct {
	Tag      string
	Text     string
	Children []*Element
	Style    tcell.Style
}

type Browser struct {
	*tview.Box
	root    *Element
	content string
	scroll  int
}

func NewBrowser() *Browser {
	return &Browser{
		Box: tview.NewBox().SetBorder(true).SetTitle("r3vr Browser"),
	}
}

func (b *Browser) SetContent(content string) {
	b.content = content
	doc, _ := html.Parse(strings.NewReader(content))
	b.root = parseHTML(doc)
}

func (b *Browser) Draw(screen tcell.Screen) {
	b.Box.DrawForSubclass(screen, b)
	x, y, width, height := b.Box.GetInnerRect()
	b.drawElement(screen, b.root, x, y, width, height, b.scroll)
}

func (b *Browser) drawElement(screen tcell.Screen, el *Element, x, y, width, height, scroll int) int {
	if y-scroll >= height {
		return y
	}

	switch el.Tag {
	case "div", "p", "h1", "h2", "h3", "h4", "h5", "h6":
		y++
	}

	if el.Text != "" {
		words := strings.Fields(el.Text)
		lineStart := x
		for _, word := range words {
			if lineStart+len(word) > x+width {
				y++
				lineStart = x
			}
			if y-scroll >= 0 && y-scroll < height {
				for i, ch := range word {
					screen.SetContent(lineStart+i, y-scroll, ch, nil, el.Style)
				}
			}
			lineStart += len(word) + 1
		}
		y++
	}

	for _, child := range el.Children {
		y = b.drawElement(screen, child, x, y, width, height, scroll)
	}

	return y
}

func parseHTML(n *html.Node) *Element {
	el := &Element{
		Tag:   n.Data,
		Style: tcell.StyleDefault,
	}

	switch n.Type {
	case html.TextNode:
		el.Text = strings.TrimSpace(n.Data)
	case html.ElementNode:
		switch n.Data {
		case "h1", "h2", "h3", "h4", "h5", "h6":
			el.Style = el.Style.Bold(true)
		case "a":
			el.Style = el.Style.Foreground(tcell.ColorBlue).Underline(true)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if child := parseHTML(c); child != nil {
			el.Children = append(el.Children, child)
		}
	}

	return el
}

func fetchPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
