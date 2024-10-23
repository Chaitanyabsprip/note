// Package main provides main  
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/chaitanyabsprip/note/cmd/note/config"
	"github.com/chaitanyabsprip/note/cmd/note/views"
	"github.com/chaitanyabsprip/note/internal/note"
	"github.com/chaitanyabsprip/note/internal/project"
)

const (
	version           = "v0.2.0"
	projectEnv        = "PROJECT"
	notesFileEnv      = "NOTESFILE"
	quietEnv          = "QUIET"
	editEnv           = "EDIT"
	peekHeadingsCount = "NOTES_HEADINGS_COUNT"
	peekHeadingsLevel = "NOTES_HEADINGS_LEVEL"
)

// CommandTree struct  
type CommandTree struct {
	w                 io.Writer
	getwd             func() (string, error)
	projectRepository project.Repository
	args              []string
}

// SetupCLI method  
func (cp *CommandTree) SetupCLI() (*config.Config, error) {
	c := new(config.Config)
	c.EditFile = os.Getenv(editEnv) != ""
	c.Project = os.Getenv(projectEnv)
	c.Quiet = os.Getenv(quietEnv) != ""
	c.Notespath = os.Getenv(notesFileEnv)
	rootCmd := createRootCmd(c)
	rootCmd.AddCommand(
		createBookmarkCmd(c),
		createDumpCmd(c),
		createIssueCmd(c),
		createPeekCmd(c),
		createTodoCmd(c),
	)
	cp.makeDumpCmdDefault(rootCmd, c)
	rootCmd.SetArgs(cp.args)
	err := rootCmd.Execute()
	if err != nil {
		return nil, err
	}
	cp.determineFilepath(c)
	return c, nil
}

func (cp *CommandTree) makeDumpCmdDefault(rootCmd *cobra.Command, c *config.Config) {
	if !c.EditFile && (len(cp.args) == 0 || cp.args[0] == "help" || cp.args[0] == "completion" ||
		cp.args[0][0] == '_') {
		return
	}
	var cmd *cobra.Command
	var err error
	var flags []string
	if rootCmd.TraverseChildren {
		cmd, flags, err = rootCmd.Traverse(cp.args)
	} else {
		cmd, flags, err = rootCmd.Find(cp.args)
	}
	isBuiltinFlag := flagsContain(flags, "-v", "-h", "--version", "--help")
	if err != nil || cmd.Use != rootCmd.Use || isBuiltinFlag {
		return
	}
	cp.args = append([]string{createDumpCmd(c).Use}, cp.args...)
}

func flagsContain(flags []string, contains ...string) bool {
	for _, flag := range contains {
		if slices.Contains(flags, flag) {
			return true
		}
	}
	return false
}

func createRootCmd(_ *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "note",
		Short: "Make notes, todos, bookmarks, issues, right from your home.",
		Long: `
Note is a command-line tool for managing your personal notes, todos, bookmarks, and issues. It
allows you to quickly create and edit notes, keep track of your tasks, and manage your bookmarks
and issues efficiently from the command line.`,
		Example: `# Create a new note
	NOTESFILE=mynotes.md note

# Edit an existing note
	NOTESFILE=mynotes.md EDIT=1 note

# Create a new todo
	NOTESFILE=mynotes.md note todo

# Add a bookmark
	note bookmark

# Report an issue
	NOTESFILE=myissues.md note issue

# Minimise output
	QUIET=1 note`,
		Version:               version,
		Args:                  cobra.ArbitraryArgs,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

func createBookmarkCmd(c *config.Config) *cobra.Command {
	cmd := cobra.Command{
		Use:   "bookmark <url> [tag] [description]",
		Short: "Create a new bookmark",
		Long:  "Create a new bookmark to save and organize URLs or references.",
		Example: `Create a bookmark with a description
		note bookmark "" "" "OpenAI" https://www.openai.com

		Add tags to a bookmark
		note bookmark "" "ai, research" https://www.openai.com`,
		Aliases:               []string{"bm", "b"},
		Args:                  cobra.MaximumNArgs(3),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c.NoteType = note.Bookmark
			if len(args) == 0 && !c.EditFile {
				var err error
				*c, err = views.GetBookmarkConfiguration()
				return err
			}
			if len(args) > 0 {
				c.Content = args[0]
			}
			if len(args) >= 1 {
				c.Tags = args[1:]
			}
			if len(args) >= 2 {
				c.Description = args[2]
			}
			return nil
		},
	}
	return &cmd
}

func createDumpCmd(c *config.Config) *cobra.Command {
	cmd := cobra.Command{
		Use:   "dump",
		Short: "Create a new note",
		Long:  "Create a new note quickly, dumping content directly from the command line.",
		Example: `Create a new note
note dump "This is a quick note"

Create a new note and edit it
EDIT=1 note dump "This is a quick note"`,
		Aliases: []string{"d"},
		Args:    cobra.ArbitraryArgs,
		Run: func(_ *cobra.Command, args []string) {
			c.NoteType = note.Dump
			c.Content = strings.Join(args, " ")
		},
	}
	return &cmd
}

func createIssueCmd(c *config.Config) *cobra.Command {
	cmd := cobra.Command{
		Use:   "issue [title] [description] [tags]",
		Short: "Create a new issue",
		Long:  "Create a new issue to track problems, bugs, or tasks.",
		Example: `# Report a new issue with a title
note issue "Bug in login feature" "The login feature fails when..."

# Add tags to an issue
note issue "Critical bug" "This is a critical issue..." "bug,urgent"`,
		Aliases:               []string{"i"},
		Args:                  cobra.ArbitraryArgs,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c.NoteType = note.Issue
			if !c.EditFile && len(args) == 0 {
				var err error
				*c, err = views.GetIssueConfiguration()
				return err
			}
			if len(args) > 0 {
				c.Title = args[0]
			}
			if len(args) >= 1 {
				c.Content = args[1]
			}
			if len(args) >= 2 {
				c.Tags = strings.Split(args[2], ",")
			}
			return nil
		},
	}
	return &cmd
}

func createPeekCmd(c *config.Config) *cobra.Command {
	cmd := cobra.Command{
		Use:   "peek",
		Short: "Take a peek at the notes",
		Long:  "Preview your notes, todos, bookmarks, or issues without opening the files.",
		Example: `# Preview bookmarks
note peek --bookmark

# Preview issues
note p --issue

# Preview todos
note peek --todo
note peek -t
note p -t

# Preview notes
note peek --dump`,
		Aliases:   []string{"p"},
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"bookmark", "bm", "b", "issue", "i", "todo", "t", "dump", "d"},
		Run: func(_ *cobra.Command, args []string) {
			c.Peek = true
			var err error
			c.NumOfHeadings, err = strconv.Atoi(os.Getenv(peekHeadingsCount))
			if err != nil {
				c.NumOfHeadings = 3
			}
			c.Level, err = strconv.Atoi(os.Getenv(peekHeadingsLevel))
			if err != nil {
				c.Level = 2
			}
			switch args[0][0] {
			case 'b':
				c.NoteType = note.Bookmark
			case 'i':
				c.NoteType = note.Issue
			case 't':
				c.NoteType = note.Todo
			case 'd':
				c.NoteType = note.Dump
			default:
				c.NoteType = note.Dump
			}
		},
	}
	return &cmd
}

func createTodoCmd(c *config.Config) *cobra.Command {
	cmd := cobra.Command{
		Use:   "todo",
		Short: "Create a new todo item",
		Long:  "Create a new todo item to keep track of tasks and actions.",
		Example: `# Create a new todo
note todo "Finish writing documentation"

# Create a new todo and edit it
EDIT=1 note todo "Finish writing documentation"`,
		Aliases: []string{"td", "t"},
		Run: func(_ *cobra.Command, args []string) {
			c.NoteType = note.Todo
			c.Content = strings.Join(args, " ")
		},
		Args: cobra.ArbitraryArgs,
	}
	return &cmd
}

func (cp *CommandTree) determineFilepath(c *config.Config) error {
	if c.Notespath != "" {
		return nil
	}
	defaultFilename := fmt.Sprint("notes.", c.NoteType, ".md")
	defaultFilepath, err := cp.getDefaultFilepath(defaultFilename)
	if err != nil {
		log.Fatal("Could not determine working directory.")
	}
	c.Notespath = defaultFilepath
	if repoRoot := project.GetRepositoryRoot(filepath.Dir(c.Notespath)); repoRoot != "" {
		c.Notespath = filepath.Join(repoRoot, defaultFilename)
	}
	if c.Project != "" {
		project := cp.projectRepository.GetProject(c.Project)
		if project == nil {
			return errors.New("could not find the project")
		}
		c.Notespath = filepath.Join(project.Path, defaultFilename)
		return nil
	}
	return nil
}

func (cp CommandTree) getDefaultFilepath(filename string) (string, error) {
	dir, err := cp.getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filename), nil
}
