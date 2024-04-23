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
	ShowPreview bool
}

func New(content, _type, notesPath string, editFile, showPreview bool) (Note, error) {
	n := new(Note)
	n.Content = content
	n.Type = _type
	n.NotesPath = notesPath
	n.EditFile = editFile
	n.ShowPreview = showPreview
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
	var noteType noteType
	switch n.Type {
	case "bookmark":
		noteType = bookmark{}
	case "dump":
		noteType = notes{}
	case "todo":
		noteType = todo{}
	default:
		fmt.Fprintln(os.Stdout, "nothing to do")
		return nil
	}
	filepath := noteType.filepath(n.NotesPath)
	setupFile(filepath, noteType.label())
	maybeOpenEditor(n.EditFile, filepath, "nvim")
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	note, err := noteType.note(n.Content, file)
	if err != nil {
		return err
	}
	_, err = file.WriteString(note)
	if err != nil {
		return err
	}
	if !n.ShowPreview {
		render(file)
	}
	return nil
}

type noteType interface {
	filepath(string) string
	label() string
	note(string, *os.File) (string, error)
}

type bookmark struct{}

func (bookmark) filepath(notesPath string) string {
	return fmt.Sprintf("%s/notes.bookmarks.md", notesPath)
}

func (bookmark) label() string {
	return "Bookmarks"
}

func (bookmark) note(content string, file *os.File) (string, error) {
	heading, err := newHeading(file)
	if err != nil {
		return "", err
	}
	note := content
	if heading != "" {
		note = fmt.Sprintf("\n%s\n\n%s", heading, content)
	}
	return fmt.Sprintln(note), nil
}

type notes struct{}

func (notes) filepath(notesPath string) string {
	return fmt.Sprintf("%s/notes.dump.md", notesPath)
}

func (notes) label() string {
	return "Notes"
}

func (notes) note(content string, file *os.File) (string, error) {
	heading, err := newHeading(file)
	if err != nil {
		return "", err
	}
	note := wordWrap(sentenceCase(content), 80)
	if heading != "" {
		note = fmt.Sprintf("\n%s\n\n%s", heading, note)
	}
	note = fmt.Sprintln(note)
	return note, nil
}

type todo struct{}

func (todo) filepath(notesPath string) string {
	return fmt.Sprintf("%s/notes.todo.md", notesPath)
}

func (todo) label() string {
	return "Todo"
}

func (todo) note(content string, file *os.File) (string, error) {
	heading, err := newHeading(file)
	if err != nil {
		return "", err
	}

	note := wordWrap(fmt.Sprint("- [ ] ", sentenceCase(content)), 80)
	if heading != "" {
		note = fmt.Sprintf("\n%s\n\n%s", heading, note)
	}
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

func newHeading(file *os.File) (string, error) {
	content, err := preview.ReadHeadings(file, 1, 2)
	if err != nil {
		return "", err
	}
	lines := strings.Split(content, "")
	lHeading := lines[len(lines)-1]
	firstEntry := lHeading == ""
	prevTime := strings.TrimPrefix(lHeading, "## ")
	currTime := time.Now().Format("Mon, 02 Jan 2006")
	if currTime != prevTime || firstEntry {
		return fmt.Sprint("## ", currTime), nil
	}
	return "", nil
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
		sb.WriteString(". ")
	}
	return strings.TrimSpace(sb.String())
}

func render(file *os.File) error {
	content, err := preview.ReadHeadings(file, 1, 2)
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
