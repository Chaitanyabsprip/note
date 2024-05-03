package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/chaitanyabsprip/note/pkg/note"
)

func TestFlagParser(t *testing.T) {
	for _, tC := range parseArgsTestCases {
		t.Run(tC.desc, func(t *testing.T) {
			cp := ConfigurationParser{
				exit: func(int) {},
				getwd: func() (string, error) {
					return tNotespath, nil
				},
				args: tC.args,
			}
			config, err := cp.ParseArgs()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("config: %#+v\n", *config)
			t.Logf("tC.config: %#+v\n", tC.config)
			if !config.Equals(tC.config) {
				t.Fail()
			}
		})
	}
}

var (
	tNotespath         = "/path/to/note/dir"
	tAltNotespath      = "/path/to/alt/note/dir"
	parseArgsTestCases = []struct {
		desc   string
		args   []string
		config Config
	}{
		{
			"-e flag should configure EditFile: true",
			[]string{"-e"},
			withDefaults(Config{EditFile: true, Notespath: getFilepath(note.Dump)}),
		},
		{
			"-f <path> flag should configure Mode: dump",
			[]string{"-f", tAltNotespath},
			withDefaults(Config{IsDump: true, Notespath: tAltNotespath}),
		},
		{
			"t flag should configure Mode: todo",
			[]string{"t"},
			withDefaults(Config{IsTodo: true, Notespath: getFilepath(note.Todo)}),
		},
		{
			"b flag should configure Mode: bookmark",
			[]string{"b"},
			withDefaults(Config{IsBookmark: true, Notespath: getFilepath(note.Bookmark)}),
		},
		{
			"d flag should configure Mode: dump",
			[]string{"d"},
			withDefaults(Config{IsDump: true, Notespath: getFilepath(note.Dump)}),
		},
		{
			"-q flag should configure Quiet: true",
			[]string{"-q"},
			withDefaults(Config{IsDump: true, Notespath: getFilepath(note.Dump), Quiet: true}),
		},
		{
			"t -e flag should configure Mode: todo, EditFile: true",
			[]string{"t", "-e"},
			withDefaults(Config{IsTodo: true, Notespath: getFilepath(note.Todo), EditFile: true}),
		},
		{
			"b -e flag should configure Mode: bookmark, Editfile: true",
			[]string{"b", "-e"},
			withDefaults(Config{IsBookmark: true, Notespath: getFilepath(note.Bookmark), EditFile: true}),
		},
		{
			"d -e flag should configure Mode: dump, EditFile: true",
			[]string{"d", "-e"},
			withDefaults(Config{IsDump: true, Notespath: getFilepath(note.Dump), EditFile: true}),
		},
		{
			"-ef <path> flag should configure Mode: dump, Global: true, EditFile: true",
			[]string{"-ef", tAltNotespath},
			withDefaults(Config{IsDump: true, Notespath: tAltNotespath, EditFile: true}),
		},
		{
			"i should create a new issue with default Title as 'Issue'",
			[]string{"i"},
			withDefaults(Config{IsIssue: true, Notespath: getFilepath(note.Issue), Title: "Issue"}),
		},
		{
			"passing -t with string should set title",
			[]string{"i", "-t", "This is a new issue"},
			withDefaults(Config{IsIssue: true, Notespath: getFilepath(note.Issue), Title: "This is a new issue"}),
		},
		{
			"positional strings should be concatenated and set as Content",
			[]string{"i", "-t", "This is a new issue", "This is the description of the issue"},
			withDefaults(Config{
				IsIssue:   true,
				Notespath: getFilepath(note.Issue),
				Title:     "This is a new issue",
				Content:   "This is the description of the issue",
			}),
		},
	}
)

func withDefaults(config Config) Config {
	updatedConfig := config
	if updatedConfig.Notespath == "" {
		updatedConfig.Notespath = tNotespath
	}
	return updatedConfig
}

func getFilepath(mode string) string {
	return filepath.Join(tNotespath, fmt.Sprint("notes.", mode, ".md"))
}
