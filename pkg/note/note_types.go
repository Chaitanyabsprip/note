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
	var traverse func(*html.Node) (string, bool)
	traverse = func(n *html.Node) (string, bool) {
		if n.Type == html.ElementNode && n.Data == "title" {
			return n.FirstChild.Data, true
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if result, ok := traverse(c); ok {
				return result, ok
			}
		}
		return "", false
	}
	title, _ := traverse(doc)
	return title
}

type notes struct{}

func (notes) label() string {
	return "Notes"
}

func (notes) toMarkdown(content string) (string, error) {
	note := wordWrap(sentenceCase(content), wrapWidth)
	return note, nil
}

type issue struct {
	createdAt   time.Time
	title       string
	description string
	status      Status
	tags        []string
}

type Status string

const (
	Open       Status = "Open"
	Closed     Status = "Closed"
	InProgress Status = "InProgress"
)

func NewIssue(title, description string, labels []string, createdAt time.Time) *issue {
	return &issue{
		createdAt:   createdAt,
		title:       title,
		description: description,
		status:      Open,
		tags:        labels,
	}
}

func (i issue) CreatedAtFormatted() string {
	if i.createdAt.IsZero() {
		return ""
	}
	return i.createdAt.Format(time.UnixDate)
}

func (issue) label() string {
	return "Issues"
}

func (i issue) toMarkdown(content string) (string, error) {
	sb := &strings.Builder{}
	fmt.Fprint(sb, "## ", wordWrap(i.title, wrapWidth))
	fmt.Fprintln(sb, "createdAt:", i.CreatedAtFormatted())
	fmt.Fprintln(sb, "status:", i.status)
	fmt.Fprintln(sb, "labels:", strings.Join(i.tags, ", "))
	sb.WriteString("\n")
	fmt.Fprint(sb, wordWrap(content, wrapWidth))
	sb.WriteString("\n")
	sb.WriteString("### Comments")
	sb.WriteString("\n")
	sb.WriteString("---")
	return sb.String(), nil
}

type todo struct{}

func (todo) label() string {
	return "Todo"
}

func (todo) toMarkdown(content string) (string, error) {
	note := wordWrap(fmt.Sprint("- [ ] ", sentenceCase(content)), wrapWidth)
	return note, nil
}
