package main

import (
	"context"
	"fmt"
	"io"
	"log"
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
	if len(args) < 2 {
		log.Fatal("Invalid Usage")
		os.Exit(1)
	}
	cp := ConfigurationParser{
		exit:   os.Exit,
		getenv: getenv,
		getwd:  os.Getwd,
		args:   args[1:],
	}
	config, err := cp.ParseArgs()
	if err != nil {
		return 1, err
	}
	// call application
	if args[1] == "peek" {
		p := preview.New(w, config.Mode, config.Notespath, config.NumOfHeadings, config.Level)
		err = p.Peek()
	} else {
		var n note.Note
		n, err = note.New(config.Content, config.Mode, config.Notespath, config.EditFile, config.Quiet)
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
