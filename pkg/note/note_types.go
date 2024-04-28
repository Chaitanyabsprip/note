package note

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type noteType interface {
	label() string
	toMarkdown(string) (string, error)
}

type bookmark struct {
	description string
	tags        []string
}

func (bookmark) label() string {
	return "Bookmarks"
}

func (b bookmark) toMarkdown(content string) (string, error) {
	title := fetchWebpageTitle(content)
	tags := strings.Join(b.tags, ", ")
	return fmt.Sprintf("\n[%s](%s)  \ntags: %s  \n%s", title, content, tags, b.description), nil
}

func fetchWebpageTitle(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return ""
	}
	var title string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.Parent != nil && n.Parent.Data == "head" {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
	return title
}

type notes struct{}

func (notes) label() string {
	return "Notes"
}

func (notes) toMarkdown(content string) (string, error) {
	note := wordWrap(sentenceCase(content), 80)
	note = fmt.Sprintln(note)
	return note, nil
}

type todo struct{}

func (todo) label() string {
	return "Todo"
}

func (todo) toMarkdown(content string) (string, error) {
	note := wordWrap(fmt.Sprint("- [ ] ", sentenceCase(content)), 80)
	note = fmt.Sprintln(note)
	return note, nil
}
