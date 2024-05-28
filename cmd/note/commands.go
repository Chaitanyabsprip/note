// Package main provides main  
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/chaitanyabsprip/note/cmd/note/config"
	"github.com/chaitanyabsprip/note/cmd/note/views"
	"github.com/chaitanyabsprip/note/internal/note"
	"github.com/chaitanyabsprip/note/internal/project"
)

const version = "v0.1.0"

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
	cp.determinFilepath(c)
	return c, nil
}

func (cp *CommandTree) makeDumpCmdDefault(rootCmd *cobra.Command, c *config.Config) {
	if len(cp.args) == 0 || cp.args[0] == "help" || cp.args[0] == "completion" ||
		cp.args[0][0] == '_' {
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

func createRootCmd(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "note",
		Short: "Make notes, todos, bookmarks, issues, right from your home.",
		Long: `
Note is a command-line tool for managing your personal notes, todos, bookmarks, and issues. It
allows you to quickly create and edit notes, keep track of your tasks, and manage your bookmarks
and issues efficiently from the command line.`,
		Example: `# Create a new note
	note --file mynotes.txt

# Edit an existing note
	note --edit --file mynotes.md

# Create a new todo
	note todo --file mytodos.md

# Add a bookmark
	note bookmark

# Report an issue
	note issue -f myissues.md

# Minimise output
	note -q`,
		Version: version,
		Args:    cobra.ArbitraryArgs,
	}
	flags := cmd.PersistentFlags()
	flags.StringVarP(&c.Notespath, "file", "f", "", "Specify notes file")
	flags.StringVarP(&c.Project, "project", "p", "", "Specify project")
	return cmd
}

func createBookmarkCmd(c *config.Config) *cobra.Command {
	cmd := cobra.Command{
		Use:   "bookmark",
		Short: "Create a new bookmark",
		Long:  "Create a new bookmark to save and organize URLs or references.",
		Example: `Create a bookmark with a description
		note bookmark --desc "OpenAI" https://www.openai.com

		Add tags to a bookmark
		note bookmark --tags "ai, research" https://www.openai.com`,
		Aliases: []string{"bm", "b"},
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c.NoteType = note.Bookmark
			if cmd.Flags().NFlag() == 0 && len(args) == 0 {
				var err error
				c, err = views.GetBookmarkConfiguration()
				if err != nil {
					return err
				}
			}
			c.Content = strings.Join(args, " ")
			return nil
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&c.EditFile, "edit", "e", false, "Open file with $EDITOR")
	flags.BoolVarP(&c.Quiet, "quiet", "q", false, "Minimise output")
	flags.StringVarP(&c.Description, "desc", "d", "", "Description for bookmarks")
	flags.StringSliceVarP(&c.Tags, "tags", "T", []string{}, "Comma separated list of tags")
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
note dump --edit "This is a quick note"`,
		Aliases: []string{"d"},
		Args:    cobra.ArbitraryArgs,
		Run: func(_ *cobra.Command, args []string) {
			c.NoteType = note.Dump
			c.Content = strings.Join(args, " ")
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&c.EditFile, "edit", "e", false, "Open file with $EDITOR")
	flags.BoolVarP(&c.Quiet, "quiet", "q", false, "Minimise output")
	return &cmd
}

func createIssueCmd(c *config.Config) *cobra.Command {
	cmd := cobra.Command{
		Use:   "issue",
		Short: "Create a new issue",
		Long:  "Create a new issue to track problems, bugs, or tasks.",
		Example: `# Report a new issue with a title
note issue --title "Bug in login feature" "The login feature fails when..."

# Add tags to an issue
note issue --tags "bug, urgent" --title "Critical bug" "This is a critical issue..."`,
		Aliases: []string{"i"},
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c.NoteType = note.Issue
			if cmd.Flags().NFlag() == 0 && len(args) == 0 {
				var err error
				c, err = views.GetIssueConfiguration()
				if err != nil {
					return err
				}
			}
			c.Content = strings.Join(args, " ")
			return nil
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&c.EditFile, "edit", "e", false, "Open file with $EDITOR")
	flags.BoolVarP(&c.Quiet, "quiet", "q", false, "Minimise output")
	flags.StringVarP(&c.Title, "title", "t", "Issue", "Title for the issue")
	flags.StringSliceVarP(&c.Tags, "tags", "T", []string{}, "Comma separated list of tags")
	return &cmd
}

func createPeekCmd(c *config.Config) *cobra.Command {
	var (
		isBookmark bool
		isIssue    bool
		isTodo     bool
		isDump     bool
	)
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
		Aliases: []string{"p"},
		Args:    cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			if isBookmark {
				c.NoteType = note.Bookmark
			} else if isIssue {
				c.NoteType = note.Issue
			} else if isTodo {
				c.NoteType = note.Todo
			} else if isDump {
				c.NoteType = note.Dump
			} else {
				c.NoteType = note.Dump
			}
			c.Peek = true
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&isBookmark, "bookmark", "b", false, "Preview bookmarks")
	flags.BoolVarP(&isIssue, "issue", "i", false, "Preview issues")
	flags.BoolVarP(&isTodo, "todo", "t", false, "Preview todos")
	flags.BoolVarP(&isDump, "dump", "d", false, "Preview notes")
	flags.IntVarP(&c.NumOfHeadings, "headings", "n", 3, "Number of headings to preview")
	flags.IntVarP(&c.Level, "level", "l", 2, "Level of markdown heading")
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
note todo --edit "Finish writing documentation"`,
		Aliases: []string{"td", "t"},
		Run: func(_ *cobra.Command, args []string) {
			c.NoteType = note.Todo
			c.Content = strings.Join(args, " ")
		},
		Args: cobra.ArbitraryArgs,
	}
	flags := cmd.Flags()
	flags.BoolVarP(&c.Quiet, "quiet", "q", false, "Minimise output")
	flags.BoolVarP(&c.EditFile, "edit", "e", false, "Open file with $EDITOR")
	return &cmd
}

func (cp *CommandTree) determinFilepath(c *config.Config) error {
	defaultFilename := fmt.Sprint("notes.", c.NoteType, ".md")
	defaultFilepath, err := cp.getDefaultFilepath(defaultFilename)
	if err != nil {
		log.Fatal("Could not determine working directory.")
	}
	if c.Notespath == "" {
		c.Notespath = defaultFilepath
	}
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

// SetupCLI method  
//
//	func createRootCmd(c *config.Config, cp *CommandTree) *cobra.Command {
//		rootCmd := &cobra.Command{
//			Use:     "note",
//			Short:   "Make notes, todos, bookmarks, issues, right from your home.",
//			Long:    "",
//			Example: "",
//			Version: version,
//			Args:    cobra.ArbitraryArgs,
//			RunE:    cp.newNote(c),
//		}
//		flags := rootCmd.PersistentFlags()
//		flags.StringVarP(&c.Notespath, "file", "f", "", "Specify notes file")
//		flags.StringVarP(&c.Project, "project", "p", "", "Specify notes file")
//		flags = rootCmd.Flags()
//		flags.BoolVarP(&c.Quiet, "quiet", "q", false, "Minimise output")
//		flags.BoolVarP(&c.EditFile, "edit", "e", false, "Open file with $EDITOR")
//		return rootCmd
//	}
//
//	func (cp *CommandTree) newNote(c *config.Config) func(_ *cobra.Command, args []string) error {
//		return func(_ *cobra.Command, args []string) error {
//			err := cp.determineFilepath(c)
//			if err != nil {
//				return nil
//			}
//			c.Content = strings.Join(args, " ")
//			var n note.Note
//			n, err = note.New(
//				c.Content,
//				c.Description,
//				c.Notespath,
//				c.Title,
//				c.NoteType(),
//				c.Tags,
//				c.EditFile,
//				c.Quiet,
//			)
//			if err != nil {
//				return err
//			}
//			err = n.Note()
//			if err != nil {
//				return err
//			}
//			return nil
//		}
//	}
//
//	func (cp CommandTree) determineFilepath(config *config.Config) error {
//		defaultFilename := fmt.Sprint("notes.", config.NoteType(), ".md")
//		defaultFilepath, err := cp.getDefaultFilepath(defaultFilename)
//		if err != nil {
//			return err
//		}
//		if config.Notespath == "" {
//			config.Notespath = defaultFilepath
//		}
//		if repoRoot := project.GetRepositoryRoot(filepath.Dir(config.Notespath)); repoRoot != "" {
//			config.Notespath = filepath.Join(repoRoot, defaultFilename)
//		}
//		if config.Project != "" {
//			project := cp.projectRepository.GetProject(config.Project)
//			if project == nil {
//				return errors.New("could not find the project")
//			}
//			config.Notespath = filepath.Join(project.Path, defaultFilename)
//			return nil
//		}
//		name := filepath.Base(filepath.Dir(config.Notespath))
//		if _, err = cp.projectRepository.AddProject(name, filepath.Dir(config.Notespath), ""); err != nil &&
//			!project.AlreadyExists(err) {
//			return err
//		}
//		return nil
//	}
