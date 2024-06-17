package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/chaitanyabsprip/note/cmd/note/config"
	"github.com/chaitanyabsprip/note/internal/note"
	"github.com/chaitanyabsprip/note/internal/project"
)

var (
	tNotespath         = "/path/to/note/dir"
	tAltNotespath      = "/path/to/alt/note/dir"
	parseArgsTestCases = []struct {
		desc   string
		args   []string
		config config.Config
	}{
		{
			"with empty arguments, config should be empty (with defaults)",
			[]string{},
			config.Config{Notespath: getFilepath(""), Title: "Issue", NumOfHeadings: 3, Level: 2},
		},
		{
			"with reserved help argument, config should be empty (with defaults)",
			[]string{"help"},
			config.Config{Notespath: getFilepath(""), Title: "Issue", NumOfHeadings: 3, Level: 2},
		},
		{
			"with reserved help argument, config should be empty (with defaults)",
			[]string{"completion"},
			config.Config{Notespath: getFilepath(""), Title: "Issue", NumOfHeadings: 3, Level: 2},
		},
		{
			"with non-flag arguments, content should be set to concatenated arguments with defaults",
			[]string{"hello"},
			withDefaults(config.Config{Content: "hello"}),
		},
		{
			"with '-e' flag, EditFile should be true, along with defaults",
			[]string{"-e"},
			withDefaults(
				config.Config{NoteType: note.Dump, Notespath: getFilepath("dump"), EditFile: true},
			),
		},
		{
			"with '-f' flag, path should be set to the argument passed with the flag, tNotespath",
			[]string{"-f", tNotespath},
			withDefaults(config.Config{NoteType: note.Dump, Notespath: tNotespath}),
		},
		{
			"with '-ef' flag, path should be set to the argument, tNotespath and EditFile as true",
			[]string{"-ef", tNotespath},
			withDefaults(config.Config{NoteType: note.Dump, Notespath: tNotespath, EditFile: true}),
		},
		{
			"with '-f' flag, path should be set to the argument passed with the flag, tAltNotespath",
			[]string{"-f", tAltNotespath},
			withDefaults(config.Config{NoteType: note.Dump, Notespath: tAltNotespath}),
		},
		{
			"with '-ef' flag, path should be set to the argument, tNotespath and EditFile as true",
			[]string{"-ef", tAltNotespath},
			withDefaults(
				config.Config{NoteType: note.Dump, Notespath: tAltNotespath, EditFile: true},
			),
		},
		{
			"with '-q' flag, Quiet should be true, along with defaults",
			[]string{"-q"},
			withDefaults(config.Config{Quiet: true}),
		},
		{
			"with '-qf' flag, path should be set to the argument, tNotespath and Quiet as true",
			[]string{"-qf", tNotespath},
			withDefaults(
				config.Config{NoteType: note.Dump, Notespath: tNotespath, Quiet: true},
			),
		},
		{
			"with trailing non-flag arguments, Content must be concatenated arguments",
			[]string{"-q", "This", "is", "content"},
			withDefaults(config.Config{Quiet: true, Content: "This is content"}),
		},
		{
			"with '-qf' flag, path should be set to the argument, tNotespath and Quiet as true",
			[]string{"-qf", tAltNotespath},
			withDefaults(
				config.Config{NoteType: note.Dump, Notespath: tAltNotespath, Quiet: true},
			),
		},
		{
			"with '-qf' flag, path should be set to the argument, tNotespath and Quiet as true",
			[]string{"-qf", tAltNotespath},
			withDefaults(
				config.Config{NoteType: note.Dump, Notespath: tAltNotespath, Quiet: true},
			),
		},
		{
			"with dump subcommand, Notespath should be <pwd>/notes.dump.md, IsDump should be true",
			[]string{"dump"},
			withDefaults(config.Config{}),
		},
		{
			"with dump subcommand and args, content should be concatenated args with defaults",
			[]string{"dump", "hello", "how"},
			withDefaults(config.Config{Content: "hello how"}),
		},
		{
			"with dump subcommand and '-q' flag, Quiet should be true with dump noteType",
			[]string{"dump", "-q"},
			withDefaults(
				config.Config{
					NoteType:  note.Dump,
					Notespath: getFilepath("dump"),
					Quiet:     true,
				},
			),
		},
		{
			"with dump subcommand, trailing non-flag arguments, Content must be concatenated arguments",
			[]string{"dump", "-q", "This", "is", "content"},
			withDefaults(config.Config{
				NoteType:  note.Dump,
				Notespath: getFilepath("dump"),
				Quiet:     true,
				Content:   "This is content",
			}),
		},
		{
			"with dump subcommand and '-e' flag, EditFile should be true with dump noteType",
			[]string{"dump", "-e"},
			withDefaults(
				config.Config{
					NoteType:  note.Dump,
					Notespath: getFilepath("dump"),
					EditFile:  true,
				},
			),
		},
		{
			"with dump subcommand and '-f' flag, path should be set to the argument passed, tNotespath",
			[]string{"dump", "-f", tNotespath},
			withDefaults(config.Config{NoteType: note.Dump, Notespath: tNotespath}),
		},
		{
			"with peek subcommand, Notespath should be <pwd>/notes.dump.md, Peek should be true",
			[]string{"peek"},
			withDefaults(config.Config{Peek: true}),
		},
		{
			"with peek subcommand and '-b' flag, Notespath should be <pwd>/notes.bookmark.md, Peek should be true",
			[]string{"peek", "-b"},
			withDefaults(
				config.Config{
					NoteType:  note.Bookmark,
					Peek:      true,
					Notespath: getFilepath("bookmark"),
				},
			),
		},
		{
			"with peek subcommand and '-i' flag, Notespath should be <pwd>/notes.issue.md, Peek should be true",
			[]string{"peek", "-i"},
			withDefaults(
				config.Config{NoteType: note.Issue, Peek: true, Notespath: getFilepath("issue")},
			),
		},
		{
			"with peek subcommand and '-t' flag, Notespath should be <pwd>/notes.todo.md, Peek should be true",
			[]string{"peek", "-t"},
			withDefaults(
				config.Config{NoteType: note.Todo, Peek: true, Notespath: getFilepath("todo")},
			),
		},
		{
			"with peek subcommand and '-d' flag, Notespath should be <pwd>/notes.dump.md, Peek should be true",
			[]string{"peek", "-d"},
			withDefaults(
				config.Config{Peek: true},
			),
		},
		{
			"with peek subcommand and '-n' flag, NumOfHeadings should be set to arg with defaults",
			[]string{"peek", "-n", "4"},
			withDefaults(config.Config{Peek: true, NumOfHeadings: 4}),
		},
		{
			"with peek subcommand and '-l' flag, Level should be set to arg with defaults",
			[]string{"peek", "-l", "1"},
			withDefaults(config.Config{Peek: true, Level: 1}),
		},
		{
			"with peek subcommand and '-f' flag, path should be set to the argument passed with the flag, tNotespath",
			[]string{"peek", "-f", tNotespath},
			withDefaults(config.Config{Notespath: tNotespath, Peek: true}),
		},
		{
			"with todo subcommand, Notespath should be <pwd>/notes.todo.md, IsTodo should be true",
			[]string{"todo"},
			withDefaults(config.Config{NoteType: note.Todo, Notespath: getFilepath("todo")}),
		},
		{
			"with todo subcommand and args, content should be concatenated args with defaults",
			[]string{"todo", "hello", "how"},
			withDefaults(
				config.Config{
					NoteType:  note.Todo,
					Notespath: getFilepath("todo"),
					Content:   "hello how",
				},
			),
		},
		{
			"with todo subcommand and '-q' flag, Quiet should be true with todo noteType",
			[]string{"todo", "-q"},
			withDefaults(
				config.Config{
					NoteType:  note.Todo,
					Notespath: getFilepath("todo"),
					Quiet:     true,
				},
			),
		},
		{
			"with todo subcommand, trailing non-flag arguments, Content must be concatenated arguments",
			[]string{"todo", "-q", "This", "is", "content"},
			withDefaults(config.Config{
				NoteType:  note.Todo,
				Notespath: getFilepath("todo"),
				Quiet:     true,
				Content:   "This is content",
			}),
		},
		{
			"with todo subcommand and '-e' flag, EditFile should be true with todo noteType",
			[]string{"todo", "-e"},
			withDefaults(
				config.Config{
					NoteType:  note.Todo,
					Notespath: getFilepath("todo"),
					EditFile:  true,
				},
			),
		},
		{
			"with todo subcommand and '-f' flag, path should be set to the argument passed, tNotespath",
			[]string{"todo", "-f", tNotespath},
			withDefaults(config.Config{NoteType: note.Todo, Notespath: tNotespath}),
		},
		{
			"with bookmark subcommand and '-q' flag, Quiet should be true with bookmark noteType",
			[]string{"bookmark", "-q"},
			withDefaults(
				config.Config{
					NoteType:  note.Bookmark,
					Notespath: getFilepath("bookmark"),
					Quiet:     true,
				},
			),
		},
		{
			"with bookmark subcommand and args, content should be concatenated args with defaults",
			[]string{"b", "hello", "how"},
			withDefaults(
				config.Config{
					NoteType:  note.Bookmark,
					Notespath: getFilepath("bookmark"),
					Content:   "hello how",
				},
			),
		},
		{
			"with bookmark subcommand, trailing non-flag arguments, Content must be concatenated arguments",
			[]string{"bookmark", "-q", "This", "is", "content"},
			withDefaults(config.Config{
				NoteType:  note.Bookmark,
				Notespath: getFilepath("bookmark"),
				Quiet:     true,
				Content:   "This is content",
			}),
		},
		{
			"with bookmark subcommand and '-e' flag, EditFile should be true with bookmark noteType",
			[]string{"bookmark", "-e"},
			withDefaults(
				config.Config{
					NoteType:  note.Bookmark,
					Notespath: getFilepath("bookmark"),
					EditFile:  true,
				},
			),
		},
		{
			"with bookmark subcommand and '-f' flag, path should be set to the argument passed, tNotespath",
			[]string{"bookmark", "-f", tNotespath},
			withDefaults(config.Config{NoteType: note.Bookmark, Notespath: tNotespath}),
		},
		{
			"with bookmark subcommand and '-d' flag, description should be set to the argument passed",
			[]string{"bookmark", "-d", "New description"},
			withDefaults(
				config.Config{
					NoteType:    note.Bookmark,
					Notespath:   getFilepath("bookmark"),
					Description: "New description",
				},
			),
		},
		{
			"with bookmark subcommand and '-T' flag, description should be set to the argument passed",
			[]string{"bookmark", "-T", "he,ll"},
			withDefaults(
				config.Config{
					NoteType:  note.Bookmark,
					Notespath: getFilepath("bookmark"),
					Tags:      []string{"he", "ll"},
				},
			),
		},
		{
			"with issue subcommand and '-q' flag, Quiet should be true with issue noteType",
			[]string{"issue", "-q"},
			withDefaults(
				config.Config{
					NoteType:  note.Issue,
					Notespath: getFilepath("issue"),
					Quiet:     true,
				},
			),
		},
		{
			"with issue subcommand and args, content should be concatenated args with defaults",
			[]string{"i", "hello", "how"},
			withDefaults(
				config.Config{
					NoteType:  note.Issue,
					Notespath: getFilepath("issue"),
					Content:   "hello how",
				},
			),
		},
		{
			"with issue subcommand, trailing non-flag arguments, Content must be concatenated arguments",
			[]string{"issue", "-q", "This", "is", "content"},
			withDefaults(config.Config{
				NoteType:  note.Issue,
				Notespath: getFilepath("issue"),
				Quiet:     true,
				Content:   "This is content",
			}),
		},
		{
			"with issue subcommand and '-e' flag, EditFile should be true with issue noteType",
			[]string{"issue", "-e"},
			withDefaults(
				config.Config{
					NoteType:  note.Issue,
					Notespath: getFilepath("issue"),
					EditFile:  true,
				},
			),
		},
		{
			"with issue subcommand and '-f' flag, path should be set to the argument passed, tNotespath",
			[]string{"issue", "-f", tNotespath},
			withDefaults(config.Config{NoteType: note.Issue, Notespath: tNotespath}),
		},
		{
			"with issue subcommand and '-t' flag, title should be set to the argument passed",
			[]string{"issue", "-t", "New title"},
			withDefaults(
				config.Config{
					NoteType:  note.Issue,
					Notespath: getFilepath("issue"),
					Title:     "New title",
				},
			),
		},
		{
			"with issue subcommand and '-T' flag, description should be set to the argument passed",
			[]string{"issue", "-T", "he,ll"},
			withDefaults(
				config.Config{
					NoteType:  note.Issue,
					Notespath: getFilepath("issue"),
					Tags:      []string{"he", "ll"},
				},
			),
		},
	}
)

func TestFlagParser(t *testing.T) {
	for _, tC := range parseArgsTestCases {
		t.Run(tC.desc, func(t *testing.T) {
			cp := CommandTree{
				w:                 new(bytes.Buffer),
				getwd:             func() (string, error) { return tNotespath, nil },
				args:              tC.args,
				projectRepository: new(MockProjectRepository),
			}
			config, err := cp.SetupCLI()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("config    : %#+v\n", *config)
			t.Logf("tC.config : %#+v\n", tC.config)
			if !config.Equals(tC.config) {
				t.Fail()
			}
		})
	}
}

type MockProjectRepository struct{}

func (mpr *MockProjectRepository) GetProject(name string) *project.Project {
	_ = name
	return new(project.Project)
}

func (mpr *MockProjectRepository) AddProject(
	name string,
	path string,
	url string,
) (*project.Project, error) {
	_ = name
	_ = path
	_ = url
	return new(project.Project), nil
}

func (mpr *MockProjectRepository) UpdateProject(
	id int,
	name string,
	path string,
	url string,
) (*project.Project, error) {
	_ = id
	_ = name
	_ = path
	_ = url
	return new(project.Project), nil
}

func withDefaults(config config.Config) config.Config {
	updatedConfig := config
	if updatedConfig.NumOfHeadings == 0 {
		updatedConfig.NumOfHeadings = 3
	}
	if updatedConfig.Level == 0 {
		updatedConfig.Level = 2
	}
	if updatedConfig.NoteType == "" {
		updatedConfig.NoteType = note.Dump
	}
	if updatedConfig.Title == "" {
		updatedConfig.Title = "Issue"
	}
	if updatedConfig.Notespath == "" {
		updatedConfig.Notespath = getFilepath("dump")
		updatedConfig.NoteType = note.Dump
	}
	return updatedConfig
}

func getFilepath(mode string) string {
	return filepath.Join(tNotespath, fmt.Sprint("notes.", mode, ".md"))
}
