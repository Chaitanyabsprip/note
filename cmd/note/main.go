package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/chaitanyabsprip/note/pkg/note"
	"github.com/chaitanyabsprip/note/pkg/preview"
)

func main() {
	exitCode, err := run(context.Background(), os.Args, os.Stdout, os.Getenv)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCode)
	}
}

func run(ctx context.Context, args []string, w io.Writer, getenv func(string) string) (int, error) {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	// parse configuration
	config, err := ParseArgs(args[1:], getenv)
	if err != nil {
		return 1, err
	}
	// call application
	if args[1] == "peek" {
		p := preview.New(w, config.Mode, config.Notespath, config.NumOfHeadings, config.Level)
		err = p.Peek()
	} else {
		n, err := note.New(config.Content, config.Mode, config.Notespath, config.EditFile, config.Quiet)
		if err != nil {
			return 1, err
		}
		err = n.Note()
		if err != nil {
			return 1, err
		}
	}
	if err != nil {
		return 1, err
	}
	return 0, nil
}
