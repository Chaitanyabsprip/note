package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/chaitanyabsprip/note/pkg/notes"
)

func main() {
	exitCode, err := run(context.Background(), os.Args, os.Stdout, os.Getenv)
	if err != nil {
		fmt.Println("There was an error", err.Error())
		os.Exit(exitCode)
	}
}

func run(ctx context.Context, args []string, w io.Writer, getenv func(string) string) (int, error) {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	config, err := parseArgs(flags, args[1:], getenv)
	if err != nil {
		return 1, err
	}

	if config.Help {
		flags.Usage()
		os.Exit(0)
	}

	notes.App(w, *config)
	return 0, nil
}

func parseArgs(flags *flag.FlagSet, args []string, getenv func(string) string) (*notes.Config, error) {
	notesfile := getDefaultNotesFile(getenv)
	config := new(notes.Config)
	flags.BoolVar(&config.Help, "h", false, "Show this help message")
	flags.BoolVar(&config.Help, "help", false, "Show this help message")
	flags.BoolVar(&config.OpenEditor, "e", false, "Open notes file in $EDITOR")
	flags.BoolVar(&config.OpenEditor, "edit", false, "Open notes file in $EDITOR")
	flags.IntVar(&config.NumOfHeadings, "n", 3, "Number of trailing headings to dump")
	flags.IntVar(&config.Level, "l", 2, "Level of heading to match (default: 2, i.e. ##)")
	flags.IntVar(&config.Level, "level", 2, "Level of heading to match (default: 2, i.e. ##)")
	flags.StringVar(&config.Filepath, "f", notesfile, "Path to the markdown file to dump headings from")
	flags.StringVar(&config.Filepath, "file", notesfile, "Path to the markdown file to dump headings from")
	flags.Parse(args)
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getDefaultNotesFile(getenv func(string) string) string {
	notesfile := getenv("NOTESFILE")
	if notesfile == "" {
		notespath := getenv("NOTESPATH")
		if notespath != "" {
			notesfile = filepath.Join(getenv("NOTESPATH"), "notes.dump.md")
		}
	}
	return notesfile
}
