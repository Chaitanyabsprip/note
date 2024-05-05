package note

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	// Bookmark  
	Bookmark = "bookmark"
	// Dump  
	Dump = "dump"
	// Issue  
	Issue = "issue"
	// Todo  
	Todo = "todo"
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
	if title == "" {
		title = content
	}
	tags := strings.Join(b.tags, ", ")
	tagsLine := "tags:"
	if tags != "" {
		tagsLine = fmt.Sprintf("tags: %s  \n", tags)
	}
	return fmt.Sprintf(
		"\n[%s](%s)  \n%s%s",
		title,
		content,
		tagsLine,
		b.description,
	), nil
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

// Status  
type Status string

const (
	// Open  
	Open Status = "Open"
	// Closed  
	Closed Status = "Closed"
	// InProgress  
	InProgress Status = "InProgress"
)

// newIssue function  
func newIssue(
	title, description string,
	labels []string,
	now time.Time,
) *issue {
	return &issue{
		createdAt:   now,
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
	fmt.Fprintln(sb, "\n##", wordWrap(sentenceCase(i.title), wrapWidth))
	fmt.Fprintln(sb, "\ncreatedAt:", i.CreatedAtFormatted())
	fmt.Fprintln(sb, "status:", i.status)
	fmt.Fprintln(sb, "labels:", strings.Join(i.tags, ", "))
	sb.WriteString("\n")
	fmt.Fprint(sb, wordWrap(content, wrapWidth))
	sb.WriteString("\n\n")
	sb.WriteString("### Comments")
	sb.WriteString("\n")
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
