package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/chaitanyabsprip/note/internal/note"
	"github.com/chaitanyabsprip/note/internal/preview"
	"github.com/chaitanyabsprip/note/internal/project"
)

func main() {
	exitCode, err := run(context.Background(), os.Args[1:], os.Getwd, os.Stdout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCode)
	}
}

func run(
	ctx context.Context,
	args []string,
	getwd func() (string, error),
	stdout io.Writer,
) (int, error) {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	cachefile, err := getConfigFilepath()
	if err != nil {
		return 1, err
	}
	pr, err := project.NewProjectRepository(cachefile)
	if err != nil {
		return 1, err
	}
	cp := CommandTree{getwd: getwd, w: stdout, args: args, projectRepository: pr}
	c, err := cp.SetupCLI()
	if err != nil {
		return 1, err
	}
	name := filepath.Base(filepath.Dir(c.Notespath))
	if _, err = cp.projectRepository.AddProject(name, filepath.Dir(c.Notespath), ""); err != nil &&
		!project.AlreadyExists(err) {
		return 1, err
	}

	if c.Peek {
		p := preview.New(
			cp.w,
			c.NoteType,
			c.Notespath,
			c.NumOfHeadings,
			c.Level,
		)
		err = p.Peek()
		if err != nil {
			return 1, err
		}
		return 0, nil
	}

	var n note.Note
	n, err = note.New(
		c.Content,
		c.Description,
		c.Notespath,
		c.Title,
		c.NoteType,
		c.Tags,
		c.EditFile,
		c.Quiet,
	)
	if err != nil {
		return 1, err
	}
	err = n.Note()
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
