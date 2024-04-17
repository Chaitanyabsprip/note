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

func run(ctx context.Context, args []string, w io.Writer, getenv func(string) string) (int, error) {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	notesfile := getenv("NOTESFILE")
	if notesfile == "" {
		notespath := getenv("NOTESPATH")
		if notespath != "" {
			notesfile = filepath.Join(getenv("NOTESPATH"), "notes.dump.md")
		}
	}
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	filepath := flags.String("f", notesfile, "Path to the markdown file to dump headings from")
	num := flags.Int("n", 3, "Number of trailing headings to dump")
	level := flags.Int("l", 2, "Level of heading to match (default: 2, e.g. ##)")
	help := flags.Bool("h", false, "Show help message")
	openEditor := flags.Bool("e", false, "Open notes file $EDITOR")

	flags.Parse(args[1:])

	if *help {
		flags.Usage()
		os.Exit(0)
	}

	config, err := notes.NewConfig(*filepath, *level, *num, *help, *openEditor)
	if err != nil {
		return 1, err
	}
	notes.App(w, *config)
	return 0, nil
}

func main() {
	exitCode, err := run(context.Background(), os.Args, os.Stdout, os.Getenv)
	if err != nil {
		fmt.Println("There was an error", err.Error())
		os.Exit(exitCode)
	}
}
