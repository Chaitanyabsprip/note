package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rsc.io/getopt"
)

type ConfigurationParser struct {
	exit   func(int)
	getenv func(string) string
	getwd  func() (string, error)
	args   []string
}

func (cp ConfigurationParser) ParseArgs() (*Config, error) {
	var config *Config
	var err error
	root := getopt.NewFlagSet("note", flag.ContinueOnError)
	config, err = parseRootArgs(root, cp.args)
	if err != nil {
		return nil, err
	}
	pArgs := root.Args()
	if root.NArg() > 0 && pArgs[0] == "peek" {
		// if peek is the first word of the note then it needs to be quoted with
		// other strings.
		previewCmd := getopt.NewFlagSet("note peek", flag.ContinueOnError)
		config, err = parsePreviewArgs(previewCmd, root.Args()[1:])
		if err != nil {
			return nil, err
		}
	}
	cp.determineFilepath(config, cp.getenv)
	return config, err
}

func parseRootArgs(flags *getopt.FlagSet, args []string) (*Config, error) {
	config := new(Config)
	config.Mode = "dump"
	addHelpFlags(flags)
	bookmark := flags.Bool("bookmark", false, "Add new bookmark")
	flags.Alias("b", "bookmark")
	dump := flags.Bool("dump", false, "Add new dump")
	flags.Alias("d", "dump")
	todo := flags.Bool("todo", false, "Add new todo")
	flags.Alias("t", "todo")
	flags.BoolVar(&config.Global, "g", false, "Use global notes")
	flags.BoolVar(&config.Quiet, "quiet", false, "Minimise output")
	flags.Alias("q", "quiet")
	flags.BoolVar(&config.EditFile, "edit", false, "Open file with $EDITOR")
	flags.Alias("e", "edit")
	flags.StringVar(&config.Notespath, "file", "", "Specify notes file")
	flags.Alias("f", "file")
	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}
	config.Mode = getNoteType(*bookmark, *dump, *todo)
	config.Content = strings.Join(flags.Args(), " ")
	return config, nil
}

func parsePreviewArgs(flags *getopt.FlagSet, args []string) (*Config, error) {
	config := new(Config)
	config.Mode = "dump"
	addHelpFlags(flags)
	bookmark := flags.Bool("bookmark", false, "Add new bookmark")
	flags.Alias("b", "bookmark")
	dump := flags.Bool("dump", false, "Add new dump")
	flags.Alias("d", "dump")
	todo := flags.Bool("todo", false, "Add new todo")
	flags.Alias("t", "todo")
	flags.BoolVar(&config.Global, "g", false, "Use global notes")
	flags.IntVar(&config.Level, "level", 2, "Level of markdown heading")
	flags.Alias("l", "level")
	flags.IntVar(&config.NumOfHeadings, "n", 3, "Number of headings to preview")
	flags.StringVar(&config.Notespath, "file", "", "Specify notes file")
	flags.Alias("f", "file")
	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}
	config.Mode = getNoteType(*bookmark, *dump, *todo)
	return config, err
}

func getNoteType(bookmark, dump, todo bool) string {
	if bookmark {
		return "bookmark"
	} else if dump {
		return "dump"
	} else if todo {
		return "todo"
	}
	return "dump"
}

func (cp ConfigurationParser) determineFilepath(config *Config, getenv func(string) string) error {
	defaultFilename := fmt.Sprint("notes.", config.Mode, ".md")
	defaultFilepath, err := cp.getDefaultFilepath(defaultFilename)
	if err != nil {
		return err
	}
	if config.Notespath == "" {
		config.Notespath = defaultFilepath
	}
	if config.Global {
		config.Notespath = filepath.Join(getenv("NOTESPATH"), defaultFilename)
	}
	return nil
}

func addHelpFlags(flags *getopt.FlagSet) {
	flags.BoolFunc("help", "Show this help message", func(s string) error {
		if s == "true" {
			flags.Usage()
			os.Exit(0)
		}
		return nil
	})
	flags.Alias("h", "help")
}

func (cp ConfigurationParser) getDefaultFilepath(filename string) (string, error) {
	dir, err := cp.getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filename), nil
}
