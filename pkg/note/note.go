package note

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chaitanyabsprip/note/pkg/notes"
)

func App(config Config) (err error) {
	switch config.Mode {
	case "bookmark":
		err = bookmark(config)
	case "dump":
		err = dump(config)
	case "todo":
		err = todo(config)
	default:
		fmt.Fprintln(os.Stdout, "nothing to do")
	}
	return err
}

func isFileClosed(file *os.File) bool {
	_, err := file.Stat()
	return err != nil
}

func bookmark(config Config) error {
	fpath := fmt.Sprintf("%s/notes.bookmarks.md", config.Notespath)
	setupFile(fpath, "Bookmarks")
	maybeOpenEditor(config, fpath, "nvim")
	file, err := os.OpenFile(fpath, os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	lHeading, err := lastHeading(file)
	if err != nil {
		return nil
	}

	firstEntry := lHeading == ""
	prevTime := strings.TrimPrefix(lHeading, "## ")
	currTime := time.Now().Format("Mon, 02 Jan 2006")

	note := config.Content
	if currTime != prevTime || firstEntry {
		note = fmt.Sprintf("\n## %s\n\n%s", currTime, config.Content)
	}
	note = fmt.Sprintln(note)

	_, err = file.WriteString(note)
	if err != nil {
		return err
	}
	if !config.Quiet {
		preview(file)
	}
	return nil
}

func dump(config Config) error {
	fpath := fmt.Sprintf("%s/notes.dump.md", config.Notespath)
	setupFile(fpath, "Notes")
	maybeOpenEditor(config, fpath, "nvim")
	file, err := os.OpenFile(fpath, os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	lHeading, err := lastHeading(file)
	if err != nil {
		return err
	}
	isFirstEntry := lHeading == ""
	prevTime := strings.TrimPrefix(lHeading, "## ")
	currTime := time.Now().Format("Mon, 02 Jan 2006")
	note := wordWrap(sentenceCase(config.Content), 80)
	if currTime != prevTime || isFirstEntry {
		note = fmt.Sprintf("\n## %s\n\n%s", currTime, note)
	}
	note = fmt.Sprintln(note)
	br, err := file.WriteString(note)
	if br < len(note) || err != nil {
		return err
	}
	if !config.Quiet {
		preview(file)
	}
	return nil
}

func todo(config Config) error {
	fpath := fmt.Sprintf("%s/notes.todo.md", config.Notespath)
	setupFile(fpath, "To-do")
	maybeOpenEditor(config, fpath, "nvim")
	file, err := os.OpenFile(fpath, os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	lHeading, err := lastHeading(file)
	if err != nil {
		return err
	}

	firstEntry := lHeading == ""
	prevTime := strings.TrimPrefix(lHeading, "## ")
	currTime := time.Now().Format("Mon, 02 Jan 2006")

	note := wordWrap(fmt.Sprint("- [ ] ", sentenceCase(config.Content)), 80)
	if currTime != prevTime || firstEntry {
		note = fmt.Sprintf("\n## %s\n\n%s", currTime, note)
	}
	note = fmt.Sprintln(note)

	_, err = file.WriteString(note)
	if err != nil {
		return err
	}
	if !config.Quiet {
		preview(file)
	}
	return nil
}

func setupFile(fpath, label string) {
	dpath := filepath.Dir(fpath)
	if _, err := os.Stat(dpath); os.IsNotExist(err) {
		os.MkdirAll(dpath, 0o755)
	}
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		heading := fmt.Sprintf("# %s\n", sentenceCase(label))
		os.WriteFile(fpath, []byte(heading), 0o644)
	}
}

func maybeOpenEditor(config Config, fpath, editorCommand string) {
	if !config.EditFile {
		return
	}
	cmd := exec.Command(editorCommand, fpath)
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

func lastHeading(file *os.File) (string, error) {
	content, err := notes.ReadHeadings(file, 1, 2)
	if err != nil {
		return "", err
	}
	lines := strings.Split(content, "")
	return lines[len(lines)-1], nil
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

func preview(file *os.File) error {
	content, err := notes.ReadHeadings(file, 1, 2)
	if err != nil {
		return err
	}
	notes.Preview(os.Stdout, content)
	return nil
}
