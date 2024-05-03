package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rsc.io/getopt"

	"github.com/chaitanyabsprip/note/pkg/project"
)

type ConfigurationParser struct {
	exit              func(int)
	getenv            func(string) string
	getwd             func() (string, error)
	projectRepository *project.ProjectRepository
	args              []string
}

func (cp ConfigurationParser) ParseArgs() (*Config, error) {
	config := new(Config)
	rootFlags := getopt.NewFlagSet("note", flag.ContinueOnError)
	registerRootFlags(rootFlags, config)
	err := rootFlags.Parse(cp.args)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	config.Content = strings.Join(rootFlags.Args(), " ")
	if rootFlags.NArg() > 0 {
		arg := rootFlags.Arg(0)
		switch arg {
		// if any subcommand is the first word of the note then it needs to be
		// quoted with other strings.
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
			registerBookmarkFlags(cmd, config)
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
		case "i", "issue":
			cmd := getopt.NewFlagSet("note issue", flag.ContinueOnError)
			registerRootFlags(cmd, config)
			registerIssueFlags(cmd, config)
			err = cmd.Parse(rootFlags.Args()[1:])
			if err != nil {
				return nil, err
			}
			config.IsIssue = true
			config.Content = strings.Join(cmd.Args(), " ")
		}
	}
	cp.determineFilepath(config, cp.getenv)
	return config, err
}

func registerIssueFlags(flags *getopt.FlagSet, config *Config) {
	flags.StringVar(&config.Title, "title", "Issue", "Title for the issue")
	flags.Alias("t", "title")
	flags.Func("tags", "Comma separated list of tags", func(s string) error {
		config.Tags = append(config.Tags, strings.Split(s, ",")...)
		return nil
	})
	flags.Alias("T", "tags")
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
	flags.BoolVar(&config.Global, "g", false, "Use global notes")
	flags.BoolVar(&config.Quiet, "quiet", false, "Minimise output")
	flags.Alias("q", "quiet")
	flags.BoolVar(&config.EditFile, "edit", false, "Open file with $EDITOR")
	flags.Alias("e", "edit")
	flags.StringVar(&config.Notespath, "file", "", "Specify notes file")
	flags.Alias("f", "file")
	flags.StringVar(&config.Project, "project", "", "Specify notes file")
	flags.Alias("p", "project")
}

func registerPreviewFlags(flags *getopt.FlagSet, config *Config) {
	addHelpFlags(flags)
	flags.BoolVar(&config.IsBookmark, "bookmark", false, "Add new bookmark")
	flags.Alias("b", "bookmark")
	flags.BoolVar(&config.IsDump, "dump", false, "Add new dump")
	flags.Alias("d", "dump")
	flags.BoolVar(&config.IsTodo, "todo", false, "Add new todo")
	flags.Alias("t", "todo")
	flags.BoolVar(&config.IsIssue, "issue", false, "Add new issue")
	flags.Alias("i", "issue")
	flags.BoolVar(&config.Global, "g", false, "Use global notes")
	flags.IntVar(&config.Level, "level", 2, "Level of markdown heading")
	flags.Alias("l", "level")
	flags.IntVar(&config.NumOfHeadings, "n", 3, "Number of headings to preview")
	flags.StringVar(&config.Notespath, "file", "", "Specify notes file")
	flags.Alias("f", "file")
	flags.StringVar(&config.Project, "project", "", "Specify notes file")
	flags.Alias("p", "project")
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
	if config.Project != "" {
		project := cp.projectRepository.GetProject(config.Project)
		if project == nil {
			return errors.New("could not find the project")
		}
		config.Notespath = filepath.Join(project.Path, defaultFilename)
		return nil
	}
	name := filepath.Base(filepath.Dir(config.Notespath))
	fmt.Println(config.Notespath)
	_, err = cp.projectRepository.AddProject(name, filepath.Dir(config.Notespath), "")
	if err != nil {
		fmt.Println(err.Error())
		return err
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
