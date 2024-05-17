// Package main provides main  
package main

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/chaitanyabsprip/note/cmd/note/config"
	"github.com/chaitanyabsprip/note/cmd/note/views"
	"github.com/chaitanyabsprip/note/pkg/note"
	"github.com/chaitanyabsprip/note/pkg/preview"
	"github.com/chaitanyabsprip/note/pkg/project"
)

// CommandTree struct  
type CommandTree struct {
	w                 io.Writer
	getwd             func() (string, error)
	projectRepository *project.Repository
}

// SetupCLI method  
func (cp *CommandTree) SetupCLI(c *config.Config) (*cobra.Command, error) {
	rootCmd := createRootCmd(c, cp)
	createBookmarkCmd(c, cp, rootCmd)
	createDumpCmd(c, cp, rootCmd)
	createTodoCmd(c, cp, rootCmd)
	createIssueCmd(c, cp, rootCmd)
	createPeekCmd(c, cp, rootCmd)
	return rootCmd, nil
}

func createRootCmd(c *config.Config, cp *CommandTree) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "note",
		Short:   "Make notes, todos, bookmarks, issues, right from your home.",
		Long:    "",
		Example: "",
		Version: "v0.1.0",
		Args:    cobra.ArbitraryArgs,
		RunE:    cp.newNote(c),
	}
	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&c.Notespath, "file", "f", "", "Specify notes file")
	flags.StringVarP(&c.Project, "project", "p", "", "Specify notes file")
	flags = rootCmd.Flags()
	flags.BoolVarP(&c.Quiet, "quiet", "q", false, "Minimise output")
	flags.BoolVarP(&c.EditFile, "edit", "e", false, "Open file with $EDITOR")
	return rootCmd
}

func createBookmarkCmd(c *config.Config, cp *CommandTree, rootCmd *cobra.Command) {
	bookmarkCmd := &cobra.Command{
		Use:     "bookmark",
		Short:   "Create a new bookmark",
		Long:    "",
		Example: "",
		Aliases: []string{"bm", "b"},
		PreRunE: func(_ *cobra.Command, args []string) error {
			var err error
			fmt.Println(args)
			c.IsBookmark = true
			if len(args) == 0 {
				c, err = views.GetBookmarkConfiguration()
				if err != nil {
					return err
				}
			}
			return nil
		},
		RunE: cp.newNote(c),
	}
	flags := bookmarkCmd.Flags()
	rootCmd.AddCommand(bookmarkCmd)
	flags.BoolVarP(&c.Quiet, "quiet", "q", false, "Minimise output")
	flags.BoolVarP(&c.EditFile, "edit", "e", false, "Open file with $EDITOR")
	flags.StringVarP(&c.Description, "desc", "d", "", "Description for bookmarks")
	flags.StringSliceVarP(&c.Tags, "tags", "T", []string{}, "Comma separated list of tags")
}

func createDumpCmd(c *config.Config, cp *CommandTree, rootCmd *cobra.Command) {
	dumpCmd := &cobra.Command{
		Use:     "dump",
		Short:   "Create a new note",
		Long:    "",
		Example: "",
		Aliases: []string{"d"},
		Args:    cobra.ArbitraryArgs,
		PreRun:  func(_ *cobra.Command, _ []string) { c.IsDump = true },
		RunE:    cp.newNote(c),
	}
	flags := dumpCmd.Flags()
	rootCmd.AddCommand(dumpCmd)
	flags.BoolVarP(&c.Quiet, "quiet", "q", false, "Minimise output")
	flags.BoolVarP(&c.EditFile, "edit", "e", false, "Open file with $EDITOR")
}

func createTodoCmd(c *config.Config, cp *CommandTree, rootCmd *cobra.Command) {
	todoCmd := &cobra.Command{
		Use:     "todo",
		Short:   "Create a new todo",
		Long:    "",
		Example: "",
		Aliases: []string{"td", "t"},
		PreRun:  func(_ *cobra.Command, _ []string) { c.IsTodo = true },
		RunE:    cp.newNote(c),
		Args: func(_ *cobra.Command, _ []string) error {
			// c.Content = strings.Join(args, " ")
			return nil
		},
	}
	flags := todoCmd.Flags()
	rootCmd.AddCommand(todoCmd)
	flags.BoolVarP(&c.Quiet, "quiet", "q", false, "Minimise output")
	flags.BoolVarP(&c.EditFile, "edit", "e", false, "Open file with $EDITOR")
}

func createIssueCmd(c *config.Config, cp *CommandTree, rootCmd *cobra.Command) {
	issueCmd := &cobra.Command{
		Use:     "issue",
		Short:   "Create a new issue",
		Long:    "",
		Example: "",
		Aliases: []string{"i"},
		Args:    cobra.ArbitraryArgs,
		PreRunE: func(_ *cobra.Command, args []string) error {
			var err error
			c.IsIssue = true
			if len(args) == 0 {
				c, err = views.GetIssueConfiguration()
				if err != nil {
					return err
				}
			}
			return nil
		},
		RunE: cp.newNote(c),
	}
	flags := issueCmd.Flags()
	rootCmd.AddCommand(issueCmd)
	flags.BoolVarP(&c.Quiet, "quiet", "q", false, "Minimise output")
	flags.BoolVarP(&c.EditFile, "edit", "e", false, "Open file with $EDITOR")
	flags.StringVarP(&c.Title, "title", "t", "Issue", "Title for the issue")
	flags.StringSliceVarP(&c.Tags, "tags", "T", []string{}, "Comma separated list of tags")
}

func createPeekCmd(c *config.Config, cp *CommandTree, rootCmd *cobra.Command) {
	peekCmd := &cobra.Command{
		Use:     "peek",
		Short:   "Take a peek at the notes",
		Long:    "",
		Example: "",
		Aliases: []string{"p"},
		RunE: func(_ *cobra.Command, _ []string) error {
			p := preview.New(
				cp.w,
				c.NoteType(),
				c.Notespath,
				c.NumOfHeadings,
				c.Level,
			)
			return p.Peek()
		},
	}
	flags := peekCmd.Flags()
	rootCmd.AddCommand(peekCmd)
	flags.BoolVarP(&c.IsBookmark, "bookmark", "b", false, "Preview bookmarks")
	flags.BoolVarP(&c.IsDump, "dump", "d", false, "Preview notes")
	flags.BoolVarP(&c.IsTodo, "todo", "t", false, "Preview todos")
	flags.BoolVarP(&c.IsIssue, "issue", "i", false, "Preview issues")
	flags.IntVarP(&c.Level, "level", "l", 2, "Level of markdown heading")
	flags.IntVar(&c.NumOfHeadings, "n", 3, "Number of headings to preview")
}

func (cp *CommandTree) newNote(c *config.Config) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		err := cp.determineFilepath(c)
		if err != nil {
			return nil
		}
		c.Content = strings.Join(args, " ")
		var n note.Note
		n, err = note.New(
			c.Content,
			c.Description,
			c.Notespath,
			c.Title,
			c.NoteType(),
			c.Tags,
			c.EditFile,
			c.Quiet,
		)
		if err != nil {
			return err
		}
		err = n.Note()
		if err != nil {
			return err
		}
		return nil
	}
}

func (cp CommandTree) determineFilepath(config *config.Config) error {
	defaultFilename := fmt.Sprint("notes.", config.NoteType(), ".md")
	defaultFilepath, err := cp.getDefaultFilepath(defaultFilename)
	if err != nil {
		return err
	}
	if config.Notespath == "" {
		config.Notespath = defaultFilepath
	}
	if repoRoot := project.GetRepositoryRoot(filepath.Dir(config.Notespath)); repoRoot != "" {
		config.Notespath = filepath.Join(repoRoot, defaultFilename)
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
	if _, err = cp.projectRepository.AddProject(name, filepath.Dir(config.Notespath), ""); err != nil &&
		!project.AlreadyExists(err) {
		return err
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
