package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/chaitanyabsprip/note/pkg/note"
	"github.com/chaitanyabsprip/note/pkg/preview"
	"github.com/chaitanyabsprip/note/pkg/project"
)

func main() {
	exitCode, err := run(context.Background(), os.Args, os.Stdout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCode)
	}
}

func run(ctx context.Context, args []string, w io.Writer) (int, error) {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	// parse configuration
	if len(args) < 2 {
		log.Fatal("Invalid Usage")
	}
	cachefile, err := getConfigFilepath()
	if err != nil {
		return 0, err
	}
	pr, err := project.NewProjectRepository(cachefile)
	if err != nil {
		return 0, err
	}
	cp := ConfigurationParser{
		exit:              os.Exit,
		getwd:             os.Getwd,
		args:              args[1:],
		projectRepository: pr,
	}
	config, err := cp.ParseArgs()
	if err != nil {
		return 1, err
	}
	// call application
	if args[1] == "peek" || args[1] == "p" {
		p := preview.New(
			w,
			config.NoteType(),
			config.Notespath,
			config.NumOfHeadings,
			config.Level,
		)
		err = p.Peek()
	} else {
		var n note.Note
		n, err = note.New(
			config.Content,
			config.Description,
			config.Notespath,
			config.Title,
			config.NoteType(),
			config.Tags,
			config.EditFile,
			config.Quiet,
		)
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

func getConfigFilepath() (string, error) {
	configDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	noteDir := filepath.Join(configDir, "note")
	if _, err := os.Stat(noteDir); os.IsNotExist(err) {
		if err := os.Mkdir(noteDir, 0o755); err != nil {
			return "", err
		}
	}
	configFile := filepath.Join(noteDir, "projects.json")
	return configFile, nil
}
