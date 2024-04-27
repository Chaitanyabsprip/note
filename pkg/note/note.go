package note

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	path "path/filepath"
	"strings"
	"time"

	"github.com/chaitanyabsprip/note/pkg/preview"
)

type Note struct {
	Content     string
	Type        string
	NotesPath   string
	EditFile    bool
	HidePreview bool
}

func New(content, _type, notesPath string, editFile, showPreview bool) (Note, error) {
	n := new(Note)
	n.Content = content
	n.Type = _type
	n.NotesPath = notesPath
	n.EditFile = editFile
	n.HidePreview = showPreview
	err := n.validate()
	if err != nil {
		return Note{}, err
	}
	return *n, nil
}

func (n Note) validate() error {
	if n.Content == "" && !n.EditFile {
		return errors.New("nothing to note here")
	}
	return nil
}

func (n Note) Note() error {
	var note noteType
	switch n.Type {
	case "bookmark":
		note = bookmark{}
	case "dump":
		note = notes{}
	case "todo":
		note = todo{}
	default:
		fmt.Fprintln(os.Stdout, "nothing to do")
		return nil
	}
	setupFile(n.NotesPath, note.label())
	maybeOpenEditor(n.EditFile, n.NotesPath, "nvim")
	file, err := os.OpenFile(n.NotesPath, os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	markdown, err := note.toMarkdown(n.Content, file)
	if err != nil {
		return err
	}
	markdown, err = addHeading(markdown, file)
	if err != nil {
		return err
	}
	_, err = file.WriteString(markdown)
	if err != nil {
		return err
	}
	if !n.HidePreview {
		render(file)
	}
	return nil
}

type noteType interface {
	label() string
	toMarkdown(string, *os.File) (string, error)
}

type bookmark struct{}

func (bookmark) label() string {
	return "Bookmarks"
}

func (bookmark) toMarkdown(content string, file *os.File) (string, error) {
	return fmt.Sprintln("<", content, ">"), nil
}

type notes struct{}

func (notes) label() string {
	return "Notes"
}

func (notes) toMarkdown(content string, file *os.File) (string, error) {
	note := wordWrap(sentenceCase(content), 80)
	note = fmt.Sprintln(note)
	return note, nil
}

type todo struct{}

func (todo) label() string {
	return "Todo"
}

func (todo) toMarkdown(content string, file *os.File) (string, error) {
	note := wordWrap(fmt.Sprint("- [ ] ", sentenceCase(content)), 80)
	note = fmt.Sprintln(note)
	return note, nil
}

func setupFile(filepath, label string) {
	dpath := path.Dir(filepath)
	if _, err := os.Stat(dpath); os.IsNotExist(err) {
		os.MkdirAll(dpath, 0o755)
	}
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		heading := fmt.Sprintf("# %s\n", sentenceCase(label))
		os.WriteFile(filepath, []byte(heading), 0o644)
	}
}

func maybeOpenEditor(editFile bool, filepath, editorCommand string) {
	if !editFile {
		return
	}
	cmd := exec.Command(editorCommand, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if err := cmd.Process.Release(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func addHeading(body string, file *os.File) (string, error) {
	heading, err := newHeading(file)
	if err != nil {
		return "", err
	}
	note := body
	if heading != "" {
		note = fmt.Sprintf("\n%s\n\n%s", heading, note)
	}
	return note, nil
}

func newHeading(file *os.File) (string, error) {
	content, err := preview.GetHeadings(file, 1, 2)
	if err != nil {
		return "", err
	}
	lines := strings.Split(content, "\n")
	lHeading := lastHeading(lines)
	prevTime := strings.TrimPrefix(lHeading, "## ")
	currTime := time.Now().Format("Mon, 02 Jan 2006")
	if currTime != prevTime || lHeading == "" {
		return fmt.Sprint("## ", currTime), nil
	}
	return "", nil
}

func lastHeading(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "##") {
			return line
		}
	}
	return ""
}

func wordWrap(text string, lineWidth int) string {
	lines := strings.Split(text, "\n")
	wrapped := ""
	for _, line := range lines {
		words := strings.Fields(strings.TrimSpace(line))
		if len(words) == 0 {
			wrapped += line + "\n"
			continue
		}
		currLine := words[0]
		for _, word := range words[1:] {
			if len(currLine)+len(word) <= lineWidth-3 {
				currLine += " " + word
			} else {
				wrapped += currLine + "\n"
				currLine = word
			}
		}
		if currLine != "" {
			wrapped += currLine + "\n"
		}
	}
	return wrapped
}

func sentenceCase(input string) string {
	var sb strings.Builder
	sentences := strings.Split(input, ". ")
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) == 0 {
			continue
		}
		sentence = strings.ToLower(sentence)
		sb.WriteString(strings.ToUpper(string(sentence[0])))
		sb.WriteString(sentence[1:])
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func render(file *os.File) error {
	content, err := preview.GetHeadings(file, 1, 2)
	if err != nil {
		return err
	}
	preview.Render(os.Stdout, content)
	return nil
}

func getGitRoot() string {
	cmd := "git"
	args := []string{"rev-parse", "--show-toplevel"}
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
