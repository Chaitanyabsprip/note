package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/chaitanyabsprip/note/cmd/note/config"
	"github.com/chaitanyabsprip/note/pkg/project"
)

func main() {
	exitCode, err := run(context.Background(), os.Getwd, os.Stdout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCode)
	}
}

func run(ctx context.Context, getwd func() (string, error), stdout io.Writer) (int, error) {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	cachefile, err := getConfigFilepath()
	if err != nil {
		return 0, err
	}
	pr, err := project.NewProjectRepository(cachefile)
	if err != nil {
		return 1, err
	}
	cp := CommandTree{getwd: getwd, w: stdout, projectRepository: pr}
	c := new(config.Config)
	cmd, err := cp.SetupCLI(c)
	if err != nil {
		return 1, err
	}
	err = cmd.Execute()
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
