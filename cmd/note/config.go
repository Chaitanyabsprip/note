package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Content         string
	Mode            string
	Notespath       string
	defaultFilename string
	Level           int
	NumOfHeadings   int
	EditFile        bool
	Quiet           bool
}

func (c *Config) Validate() error {
	if c.Level != 0 && (c.Level > 6 || c.Level < 1) {
		return errors.New("level can be between [1, 6]")
	}

	if c.Content == "" && !c.EditFile && !c.isPreviewing() {
		return errors.New("nothing to note here")
	}
	return nil
}

func (c Config) isPreviewing() bool {
	return c.Level != 0 && c.NumOfHeadings != 0
}

func (c *Config) FilePath() string {
	c.defaultFilename = fmt.Sprint("notes.", c.Mode, ".md")
	pwd, _ := os.Getwd()
	defaultPath := filepath.Join(pwd, c.defaultFilename)
	notespath := c.Notespath
	if notespath == "" {
		notespath = defaultPath
	}
	info, err := os.Stat(notespath)
	if err != nil {
		notespath = filepath.Join(os.Getenv("NOTESPATH"), c.defaultFilename)
	}
	if !info.IsDir() {
		return notespath
	}
	notespath = filepath.Join(notespath, c.defaultFilename)
	return notespath
}
