package note

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	Bookmark = "bookmark"
	Dump     = "dump"
	Issue    = "issue"
	Todo     = "todo"
)

type noteType interface {
	Label() string
	ToMarkdown(string) (string, error)
}

type bookmark struct {
	description string
	tags        []string
}

func (bookmark) Label() string {
	return "Bookmarks"
}

func (b bookmark) ToMarkdown(content string) (string, error) {
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

func (notes) Label() string {
	return "Notes"
}

func (notes) ToMarkdown(content string) (string, error) {
	note := wordWrap(sentenceCase(content), wrapWidth)
	return note, nil
}

type issue struct {
	createdAt   time.Time
	title       string
	description string
	labels      []string
	status      Status
}

type Status int

const (
	Open Status = iota + 1
	Closed
	InProgress
)

func (issue) Label() string {
	return "Issues"
}

func (i issue) toMarkdown(content string) (string, error) {
	var sb *strings.Builder
	// ## Issue 1: Issue 1 Title
	fmt.Fprintln(sb, "##", wordWrap(i.title, wrapWidth))
	// createdAt:
	fmt.Fprintln(sb, "createdAt:", i.createdAt.Format(time.UnixDate))
	// labels:
	fmt.Fprintln(sb, "labels:", i.labels)
	//
	sb.WriteString("\n")
	// Issue 1 description goes here.
	fmt.Fprintln(sb, wordWrap(content, wrapWidth))
	//
	sb.WriteString("\n")
	// ### Comments
	sb.WriteString("### Comments")
	//
	sb.WriteString("\n")
	// ---
	sb.WriteString("---")
	return "", nil
}

type todo struct{}

func (todo) Label() string {
	return "Todo"
}

func (todo) ToMarkdown(content string) (string, error) {
	note := wordWrap(fmt.Sprint("- [ ] ", sentenceCase(content)), wrapWidth)
	return note, nil
}
