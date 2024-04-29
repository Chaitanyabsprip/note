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
	config := new(Config)
	rootFlags := getopt.NewFlagSet("note", flag.ContinueOnError)
	registerRootFlags(rootFlags, config)
	registerNoteTypeFlags(rootFlags, config)
	err := rootFlags.Parse(cp.args)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	config.Content = strings.Join(rootFlags.Args(), " ")
	if rootFlags.NArg() > 0 {
		arg := rootFlags.Arg(0)
		switch arg {
		// if peek is the first word of the note then it needs to be quoted with
		// other strings.
		case "p", "peek":
			cmd := getopt.NewFlagSet("note peek", flag.ContinueOnError)
			registerPreviewFlags(cmd, config)
			err = cmd.Parse(rootFlags.Args()[1:])
			if err != nil {
				return nil, err
			}
		case "b", "bm", "bookmark":
			cmd := getopt.NewFlagSet("note bookmark", flag.ContinueOnError)
			registerRootFlags(cmd, config)
			err = cmd.Parse(rootFlags.Args()[1:])
			if err != nil {
				return nil, err
			}
			config.IsBookmark = true
			config.Content = strings.Join(cmd.Args(), " ")
		case "d", "dump":
			cmd := getopt.NewFlagSet("note bookmark", flag.ContinueOnError)
			registerRootFlags(cmd, config)
			err = cmd.Parse(rootFlags.Args()[1:])
			if err != nil {
				return nil, err
			}
			config.IsDump = true
			config.Content = strings.Join(cmd.Args(), " ")
		case "t", "td", "todo":
			cmd := getopt.NewFlagSet("note bookmark", flag.ContinueOnError)
			registerRootFlags(cmd, config)
			err = cmd.Parse(rootFlags.Args()[1:])
			if err != nil {
				return nil, err
			}
			config.IsTodo = true
			config.Content = strings.Join(cmd.Args(), " ")
		}
	}
	cp.determineFilepath(config, cp.getenv)
	return config, err
}

func registerNoteTypeFlags(flags *getopt.FlagSet, config *Config) {
	flags.BoolVar(&config.IsBookmark, "bookmark", false, "Add new bookmark")
	flags.Alias("b", "bookmark")
	flags.BoolVar(&config.IsDump, "dump", false, "Add new dump")
	flags.Alias("d", "dump")
	flags.BoolVar(&config.IsTodo, "todo", false, "Add new todo")
	flags.Alias("t", "todo")
}

func registerBookmarkFlags(flags *getopt.FlagSet, config *Config) {
	flags.StringVar(&config.Description, "desc", "", "Description for bookmarks")
	flags.Alias("D", "desc")
	flags.Func("tags", "Comma separated list of tags", func(s string) error {
		config.Tags = append(config.Tags, strings.Split(s, ",")...)
		return nil
	})
	flags.Alias("T", "tags")
}

func registerRootFlags(flags *getopt.FlagSet, config *Config) {
	addHelpFlags(flags)
	flags.StringVar(&config.Description, "desc", "", "Description for bookmarks")
	flags.Alias("D", "desc")
	flags.Func("tags", "Comma separated list of tags", func(s string) error {
		config.Tags = append(config.Tags, strings.Split(s, ",")...)
		return nil
	})
	flags.Alias("T", "tags")
	flags.BoolVar(&config.Global, "g", false, "Use global notes")
	flags.BoolVar(&config.Quiet, "quiet", false, "Minimise output")
	flags.Alias("q", "quiet")
	flags.BoolVar(&config.EditFile, "edit", false, "Open file with $EDITOR")
	flags.Alias("e", "edit")
	flags.StringVar(&config.Notespath, "file", "", "Specify notes file")
	flags.Alias("f", "file")
}

func registerPreviewFlags(flags *getopt.FlagSet, config *Config) {
	addHelpFlags(flags)
	flags.BoolVar(&config.IsBookmark, "bookmark", false, "Add new bookmark")
	flags.Alias("b", "bookmark")
	flags.BoolVar(&config.IsDump, "dump", false, "Add new dump")
	flags.Alias("d", "dump")
	flags.BoolVar(&config.IsTodo, "todo", false, "Add new todo")
	flags.Alias("t", "todo")
	flags.BoolVar(&config.Global, "g", false, "Use global notes")
	flags.IntVar(&config.Level, "level", 2, "Level of markdown heading")
	flags.Alias("l", "level")
	flags.IntVar(&config.NumOfHeadings, "n", 3, "Number of headings to preview")
	flags.StringVar(&config.Notespath, "file", "", "Specify notes file")
	flags.Alias("f", "file")
}

func (cp ConfigurationParser) determineFilepath(config *Config, getenv func(string) string) error {
	defaultFilename := fmt.Sprint("notes.", config.NoteType(), ".md")
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
