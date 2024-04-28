package note

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	path "path/filepath"
	"strings"

	"github.com/chaitanyabsprip/note/pkg/preview"
)

type Note struct {
	Content     string
	Description string
	NotesPath   string
	Type        string
	Tags        []string
	EditFile    bool
	HidePreview bool
}

func New(content, description, notesPath, _type string, tags []string, editFile, showPreview bool) (Note, error) {
	n := new(Note)
	n.Content = content
	n.Description = description
	n.EditFile = editFile
	n.HidePreview = showPreview
	n.NotesPath = notesPath
	n.Tags = tags
	n.Type = _type
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
	note := n.getNoteType()
	if note == nil {
		return nil
	}
	setupFile(n.NotesPath, note.label())
	maybeOpenEditor(n.EditFile, n.NotesPath, "nvim")
	markdown, err := note.toMarkdown(n.Content)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(n.NotesPath, os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
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

func (n Note) getNoteType() noteType {
	var note noteType
	switch n.Type {
	case "bookmark":
		note = bookmark{
			description: n.Description,
			tags:        n.Tags,
		}
	case "dump":
		note = notes{}
	case "todo":
		note = todo{}
	default:
		fmt.Fprintln(os.Stdout, "nothing to do")
		return nil
	}
	return note
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
