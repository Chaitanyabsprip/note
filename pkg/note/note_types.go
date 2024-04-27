package note

import "fmt"

type noteType interface {
	label() string
	toMarkdown(string) (string, error)
}

type bookmark struct{}

func (bookmark) label() string {
	return "Bookmarks"
}

func (bookmark) toMarkdown(content string) (string, error) {
	return fmt.Sprint("[](", content, ")\n\n"), nil
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
