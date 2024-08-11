# r3vr: Terminal-based Web Browser

## r3vr (pronounced "revere") is a simple, terminal-based web browser written in Go. It allows users to view web pages in a text-based format directly in their terminal.

## Features

Fetch and display web pages in the terminal
Basic HTML parsing and rendering
Simple navigation through URL input
Text-based representation of web content
Basic styling (bold for headers, blue for links)

## Requirements

Go 1.16 or higher
Terminal with Unicode support

## Installation

- Clone the repository:
  Copy

  ```bash
  git clone https://github.com/yourusername/r3vr.git
  cd r3vr
  ```

- Install dependencies:
  Copy

```bash
go get github.com/gdamore/tcell/v2
go get github.com/rivo/tview
go get golang.org/x/net/html
```

- Build the project:
  Copy

```bash
go build
```

## Usage

- Run the application:
  Copy./r3vr

Enter a full URL (including http:// or https://) in the input field at the top of the screen.
Press Enter to load the page.
The content of the web page will be displayed in the main area of the terminal.

How It Works
r3vr is built using several key components:

Main Application: Uses the tview library to create the terminal user interface.
Browser Component: Handles the rendering of web content.
HTML Parsing: Utilizes Go's net/html package to parse HTML content.
Content Fetching: Uses Go's net/http package to fetch web pages.

## Detailed Workflow

User Input: The user enters a URL in the input field.
Fetching: The fetchPage function sends an HTTP GET request to the specified URL and retrieves the HTML content.
Parsing: The HTML content is parsed into a tree structure using the html.Parse function.
Element Creation: The parseHTML function recursively traverses the HTML tree, creating Element structs that represent the structure and content of the page.
Rendering: The Browser.Draw method is called to render the content. It uses the drawElement method to recursively draw each element of the page.
Styling: Basic styling is applied during rendering. Headers are made bold, and links are colored blue and underlined.

Limitations

JavaScript is not executed, so dynamic content won't be rendered.
CSS styling is not applied beyond the basic styles implemented in the code.
Images and other media are not displayed.
Complex layouts may not be accurately represented.

Contributing
Contributions to r3vr are welcome! Please feel free to submit pull requests, create issues, or suggest new features.
License
This project is licensed under the MIT License - see the LICENSE file for details.

Developer Documentation
Key Structures

Element: Represents an HTML element.

Tag: The HTML tag (e.g., "div", "p", "a")
Text: The text content of the element
Children: Slice of child elements
Style: Styling information for rendering

Browser: The main component for displaying web content.

Box: Embedded tview.Box for basic GUI functionality
root: The root Element of the parsed HTML
content: Raw HTML content
scroll: Current scroll position (not fully implemented in the current version)

## Key Functions

```go
NewBrowser(): Creates a new Browser instance.
fetchPage(url string): Fetches the HTML content of a given URL.
parseHTML(n *html.Node): Recursively parses an HTML node into our Element structure.
Browser.SetContent(content string): Sets the content of the browser and triggers parsing.
Browser.Draw(screen tcell.Screen): Draws the browser content on the screen.
Browser.drawElement(...): Recursively draws an element and its children.
```
