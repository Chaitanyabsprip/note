package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const usage = `
Usage:
    %[1]s [options]

Options:
    -b, --bookmark       Create a new bookmark
    -d, --dump           Create a new note in the dump file(default)
    -t, --todo           Create a new todo item
    -e, --edit           Edit the notes file in the default editor
    -l, --local          Use a file local to current directory
    -q, --quiet          Be silent
    -h, --help           Show this help message

Examples:
    %[1]s This is a note
    %[1]s -b 'https://newbookmark.com'
    %[1]s -t This is a new todo item
`

func ParseArgs(flags *flag.FlagSet, args []string, getenv func(string) string) (*Config, error) {
	var config *Config
	var err error
	if args[0] == "peek" {
		config, err = parsePreviewArgs(flags, args[1:], getenv)
	} else {
		config, err = parseNoteArgs(flags, args, getenv)
	}
	return config, err
}

func parsePreviewArgs(flags *flag.FlagSet, args []string, getenv func(string) string) (*Config, error) {
	config := new(Config)
	config.Mode = "dump"
	notespath := getenv("NOTESPATH")
	bookmarkFlagFunc := setMode(config, "bookmark")
	dumpFlagFunc := setMode(config, "dump")
	todoFlagFunc := setMode(config, "todo")
	flags.IntVar(&config.NumOfHeadings, "n", 3, "Number of trailing headings to dump")
	flags.IntVar(&config.Level, "l", 2, "Level of heading to match (default: 2, i.e. ##)")
	flags.IntVar(&config.Level, "level", 2, "Level of heading to match (default: 2, i.e. ##)")
	flags.StringVar(&config.Notespath, "f", notespath, "Path to the markdown file to dump headings from")
	flags.StringVar(&config.Notespath, "file", notespath, "Path to the markdown file to dump headings from")
	flags.BoolFunc("bookmark", "Add the note to the bookmark list", bookmarkFlagFunc)
	flags.BoolFunc("b", "Add the note to the bookmark list", bookmarkFlagFunc)
	flags.BoolFunc("dump", "Write the note to the notes file", dumpFlagFunc)
	flags.BoolFunc("d", "Dump the notes to a file", dumpFlagFunc)
	flags.BoolFunc("todo", "Set the note as a todo item", todoFlagFunc)
	flags.BoolFunc("t", "Set the note as a todo item", todoFlagFunc)
	flags.Parse(args)
	config.Notespath = config.FilePath()
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func parseNoteArgs(flags *flag.FlagSet, args []string, getenv func(string) string) (*Config, error) {
	config := new(Config)
	config.Mode = "dump"
	binName := filepath.Base(os.Args[0])
	notespath, err := defaultNotespath(getenv)
	if err != nil {
		return nil, err
	}
	config.Notespath = notespath
	dumpFlagFunc := setMode(config, "dump")
	bookmarkFlagFunc := setMode(config, "bookmark")
	todoFlagFunc := setMode(config, "todo")

	flags.Usage = func() { fmt.Printf(usage, binName) }
	flags.BoolVar(&config.Quiet, "quiet", false, "Be silent")
	flags.BoolVar(&config.Quiet, "q", false, "Be silent")
	flags.BoolVar(&config.EditFile, "edit", false, "Edit the notes file in the default editor")
	flags.BoolVar(&config.EditFile, "e", false, "Edit the notes file in the default editor")
	flags.StringVar(&config.Notespath, "f", notespath, "Path to the markdown file to dump headings from")
	flags.StringVar(&config.Notespath, "file", notespath, "Path to the markdown file to dump headings from")
	flags.BoolFunc("bookmark", "Add the note to the bookmark list", bookmarkFlagFunc)
	flags.BoolFunc("b", "Add the note to the bookmark list", bookmarkFlagFunc)
	flags.BoolFunc("dump", "Write the note to the notes file", dumpFlagFunc)
	flags.BoolFunc("d", "Dump the notes to a file", dumpFlagFunc)
	flags.BoolFunc("todo", "Set the note as a todo item", todoFlagFunc)
	flags.BoolFunc("t", "Set the note as a todo item", todoFlagFunc)

	localFlagFunc := func(x string) error {
		var np string
		np, err = os.Getwd()
		config.Notespath = np
		return err
	}
	flags.BoolFunc("local", "Use a file local to current directory", localFlagFunc)
	flags.BoolFunc("l", "Use a file local to current directory", localFlagFunc)
	flags.Parse(args)
	config.Content = strings.Join(flags.Args(), " ")
	config.Notespath = config.FilePath()
	err = config.Validate()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func setMode(config *Config, mode string) func(string) error {
	return func(_ string) error {
		config.Mode = mode
		return nil
	}
}

func defaultNotespath(getenv func(string) string) (notespath string, err error) {
	_ = getenv
	// notespath = getenv("NOTESPATH")
	// if notespath == "" {
	notespath, err = os.Getwd()
	// }
	return
}
