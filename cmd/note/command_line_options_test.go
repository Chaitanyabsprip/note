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
				exit:   func(int) {},
				getenv: tGetenv,
				getwd: func() (string, error) {
					return tGetenv(""), nil
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
			"-g flag should configure Global: true",
			[]string{"-g"},
			withDefaults(Config{Global: true, Notespath: getFilepath(note.Dump)}),
		},
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
			"-t flag should configure Mode: todo",
			[]string{"-t"},
			withDefaults(Config{IsTodo: true, Notespath: getFilepath(note.Todo)}),
		},
		{
			"-b flag should configure Mode: bookmark",
			[]string{"-b"},
			withDefaults(Config{IsBookmark: true, Notespath: getFilepath(note.Bookmark)}),
		},
		{
			"-d flag should configure Mode: dump",
			[]string{"-d"},
			withDefaults(Config{IsDump: true, Notespath: getFilepath(note.Dump)}),
		},
		{
			"-q flag should configure Quiet: true",
			[]string{"-q"},
			withDefaults(Config{IsDump: true, Notespath: getFilepath(note.Dump), Quiet: true}),
		},
		{
			"-ge flag should configure Global: true, EditFile: true",
			[]string{"-ge"},
			withDefaults(Config{IsDump: true, Notespath: getFilepath(note.Dump), Global: true, EditFile: true}),
		},
		{
			"-gt flag should configure Mode: todo, Global: true",
			[]string{"-gt"},
			withDefaults(Config{IsTodo: true, Notespath: getFilepath(note.Todo), Global: true}),
		},
		{
			"-gb flag should configure Mode: bookmark, Global: true",
			[]string{"-gb"},
			withDefaults(Config{IsBookmark: true, Notespath: getFilepath(note.Bookmark), Global: true}),
		},
		{
			"-gd flag should configure Mode: dump, Global: true",
			[]string{"-gd"},
			withDefaults(Config{IsDump: true, Notespath: getFilepath(note.Dump), Global: true}),
		},
		{
			"-te flag should configure Mode: todo, EditFile: true",
			[]string{"-te"},
			withDefaults(Config{IsTodo: true, Notespath: getFilepath(note.Todo), EditFile: true}),
		},
		{
			"-be flag should configure Mode: bookmark, Editfile: true",
			[]string{"-be"},
			withDefaults(Config{IsBookmark: true, Notespath: getFilepath(note.Bookmark), EditFile: true}),
		},
		{
			"-de flag should configure Mode: dump, EditFile: true",
			[]string{"-de"},
			withDefaults(Config{IsDump: true, Notespath: getFilepath(note.Dump), EditFile: true}),
		},
		{
			"-gte flag should configure Mode: todo, Global: true, EditFile: true",
			[]string{"-gte"},
			withDefaults(Config{IsTodo: true, Notespath: getFilepath(note.Todo), EditFile: true, Global: true}),
		},
		{
			"-gbe flag should configure Mode: bookmark, Global: true, EditFile: true",
			[]string{"-gbe"},
			withDefaults(Config{IsBookmark: true, Notespath: getFilepath(note.Bookmark), EditFile: true, Global: true}),
		},
		{
			"-gde flag should configure Mode: dump, Global: true, EditFile: true",
			[]string{"-gde"},
			withDefaults(Config{IsDump: true, Notespath: getFilepath(note.Dump), EditFile: true, Global: true}),
		},
		{
			"-ef <path> flag should configure Mode: dump, Global: true, EditFile: true",
			[]string{"-ef", tAltNotespath},
			withDefaults(Config{IsDump: true, Notespath: tAltNotespath, EditFile: true}),
		},
	}
)

func tGetenv(_ string) string {
	return tNotespath
}

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
