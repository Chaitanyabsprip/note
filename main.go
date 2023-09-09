package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	helpFlag  bool
	quietFlag bool
	editFlag  bool
	modeFlag  string
	notesPath string
)

const usage = `
Usage:
    %[1]s [options]

Options:
    -D, --daily          Create a daily note
    -b, --bookmark       Create a new bookmark
    -d, --dump           Create a new note in the dump file(default)
    -t, --todo           Create a new todo item
    -e, --edit           Edit the notes file in the default editor
    -l, --local          Use a file local to current directory
    -q,                  Be silent
    -h, --help           Show this help message

Examples:
    %[1]s This is a note
    %[1]s -b 'https://newbookmark.com'
    %[1]s -t This is a new todo item
`

func init() {
	modeFlag = "dump"
	home := os.Getenv("HOME")
	binName := filepath.Base(os.Args[0])
	notesPath = home + "/.cache/notes"

	setMode := func(mode string) func(string) error {
		return func(_ string) error {
			modeFlag = mode
			return nil
		}
	}

	flag.Usage = func() { fmt.Printf(usage, binName) }

	flag.BoolVar(&helpFlag, "help", false, "Show this help message")
	flag.BoolVar(&helpFlag, "h", false, "Show this help message")

	flag.BoolVar(&quietFlag, "quiet", false, "Be silent")
	flag.BoolVar(&quietFlag, "q", false, "Be silent")

	flag.BoolVar(&editFlag, "edit", false, "Edit the notes file in the default editor")
	flag.BoolVar(&editFlag, "e", false, "Edit the notes file in the default editor")

	bookmarkFlagFunc := setMode("bookmark")
	flag.BoolFunc("bookmark", "Add the note to the bookmark list", bookmarkFlagFunc)
	flag.BoolFunc("b", "Add the note to the bookmark list", bookmarkFlagFunc)

	dumpFlagFunc := setMode("dump")
	flag.BoolFunc("dump", "Write the note to the notes file", dumpFlagFunc)
	flag.BoolFunc("d", "Dump the notes to a file", dumpFlagFunc)

	todoFlagFunc := setMode("todo")
	flag.BoolFunc("todo", "Set the note as a todo item", todoFlagFunc)
	flag.BoolFunc("t", "Set the note as a todo item", todoFlagFunc)

	localFlagFunc := func(_ string) error {
		np, err := os.Getwd()
		notesPath = np
		return err
	}
	flag.BoolFunc("local", "Use a file local to current directory", localFlagFunc)
	flag.BoolFunc("l", "Use a file local to current directory", localFlagFunc)
	_ = notesPath
}

func validateContent(content string) {
	if content == "" {
		fmt.Fprintln(os.Stdout, "Nothing to note here :shrug:")
		os.Exit(0)
	}
}
func setupFile(fpath, label string) {
	dpath := filepath.Dir(fpath)
	if _, err := os.Stat(dpath); os.IsNotExist(err) {
		os.MkdirAll(dpath, 0755)
	}
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		heading := fmt.Sprintf("# %s\n", strings.ToTitle(label))
		os.WriteFile(fpath, []byte(heading), 0644)
	}
}

func maybeOpenEditor(fpath string) {
	if !editFlag {
		return
	}
	cmd := exec.Command(os.Getenv("EDITOR"), fpath)
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

func lastHeading(fpath string) (string, int, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return "", -1, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	linenr := 1
	var lHeading string

	for s.Scan() {
		line := s.Text()
		linenr++
		if strings.HasPrefix(line, "## ") {
			lHeading = line
		}
	}
	return lHeading, linenr, nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		os.Exit(1)
	}
}

func appendToFile(fpath, content string) {
	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_WRONLY, 0644)
	checkErr(err)
	defer f.Close()

	_, err = f.WriteString(content)
	checkErr(err)
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

func bookmark(content string) {
	fpath := fmt.Sprintf("%s/notes.bookmarks.md", notesPath)
	setupFile(fpath, "Bookmarks")
	maybeOpenEditor(fpath)
	validateContent(content)

	lhead, linenr, err := lastHeading(fpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	firstEntry := lhead == ""
	_ = linenr

	prevTime := strings.TrimPrefix(lhead, "## ")
	currTime := time.Now().Format("Mon, 02 Jan 2006")

	note := content
	if currTime != prevTime || firstEntry {
		note = fmt.Sprintf("\n## %s\n\n%s", currTime, content)
	}
	note = fmt.Sprintln(note)

	appendToFile(fpath, note)
}

func dump(content string) {
	fpath := fmt.Sprintf("%s/notes.dump.md", notesPath)
	setupFile(fpath, "Notes")
	maybeOpenEditor(fpath)
	validateContent(content)

	lhead, linenr, err := lastHeading(fpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	firstEntry := lhead == ""
	_ = linenr

	prevTime := strings.TrimPrefix(lhead, "## ")
	currTime := time.Now().Format("Mon, 02 Jan 2006")

	note := wordWrap(strings.ToTitle(content), 80)
	if currTime != prevTime || firstEntry {
		note = fmt.Sprintf("\n## %s\n\n%s", currTime, note)
	}
	note = fmt.Sprintln(note)

	appendToFile(fpath, note)
}

func todo(content string) {
	fpath := fmt.Sprintf("%s/notes.todo.md", notesPath)
	setupFile(fpath, "To-do")
	maybeOpenEditor(fpath)
	validateContent(content)

	lhead, linenr, err := lastHeading(fpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	firstEntry := lhead == ""
	_ = linenr

	prevTime := strings.TrimPrefix(lhead, "## ")
	currTime := time.Now().Format("Mon, 02 Jan 2006")

	sentenceCase := strings.ToUpper(string(content[0])) + string(content[1:])
	note := wordWrap(fmt.Sprint("- [ ] ", sentenceCase), 80)
	if currTime != prevTime || firstEntry {
		note = fmt.Sprintf("\n## %s\n\n%s", currTime, note)
	}
	note = fmt.Sprintln(note)

	appendToFile(fpath, note)
}

func main() {
	flag.Parse()

	if helpFlag {
		flag.Usage()
		return
	}

	content := strings.Join(flag.Args(), " ")

	switch modeFlag {
	case "bookmark":
		bookmark(content)
	case "dump":
		dump(content)
	case "todo":
		todo(content)
	default:
		fmt.Fprintln(os.Stdout, "nothing to do")
	}

}
