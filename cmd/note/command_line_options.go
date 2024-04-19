package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chaitanyabsprip/note/pkg/note"
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

func parseArgs(flags *flag.FlagSet, args []string, getenv func(string) string) (*note.Config, error) {
	config := new(note.Config)
	config.Mode = "dump"
	binName := filepath.Base(os.Args[0])
	notespath, err := defaultNotespath(getenv)
	if err != nil {
		return nil, err
	}
	config.Notespath = notespath

	setMode := func(mode string) func(string) error {
		return func(_ string) error {
			config.Mode = mode
			return nil
		}
	}

	flags.Usage = func() { fmt.Printf(usage, binName) }

	flags.BoolVar(&config.Help, "help", false, "Show this help message")
	flags.BoolVar(&config.Help, "h", false, "Show this help message")

	flags.BoolVar(&config.Quiet, "quiet", false, "Be silent")
	flags.BoolVar(&config.Quiet, "q", false, "Be silent")

	flags.BoolVar(&config.EditFile, "edit", false, "Edit the notes file in the default editor")
	flags.BoolVar(&config.EditFile, "e", false, "Edit the notes file in the default editor")

	bookmarkFlagFunc := setMode("bookmark")
	flags.BoolFunc("bookmark", "Add the note to the bookmark list", bookmarkFlagFunc)
	flags.BoolFunc("b", "Add the note to the bookmark list", bookmarkFlagFunc)

	dumpFlagFunc := setMode("dump")
	flags.BoolFunc("dump", "Write the note to the notes file", dumpFlagFunc)
	flags.BoolFunc("d", "Dump the notes to a file", dumpFlagFunc)

	todoFlagFunc := setMode("todo")
	flags.BoolFunc("todo", "Set the note as a todo item", todoFlagFunc)
	flags.BoolFunc("t", "Set the note as a todo item", todoFlagFunc)

	localFlagFunc := func(x string) error {
		fmt.Printf("%T", x)
		var np string
		np, err = os.Getwd()
		config.Notespath = np
		return err
	}
	flags.BoolFunc("local", "Use a file local to current directory", localFlagFunc)
	flags.BoolFunc("l", "Use a file local to current directory", localFlagFunc)
	flags.Parse(args)
	config.Content = strings.Join(flags.Args(), "")
	err = config.Validate()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func defaultNotespath(getenv func(string) string) (notespath string, err error) {
	notespath = getenv("NOTESPATH")
	if notespath == "" {
		notespath, err = os.Getwd()
	}
	return
}
