package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/chaitanyabsprip/note/pkg/note"
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
	_ = w
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	config, err := ParseArgs(flags, args[1:], getenv)
	if err != nil {
		return 1, err
	}
	if args[1] == "peek" {
		fmt.Println(config.Notespath)
		err = note.Peek(w, *config)
	} else {
		err = note.App(*config)
	}
	if err != nil {
		return 1, err
	}
	return 0, nil
}
