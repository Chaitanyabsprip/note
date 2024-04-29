package main

import (
	"slices"

	"github.com/chaitanyabsprip/note/pkg/note"
)

type Config struct {
	Content       string
	Description   string
	Notespath     string
	Tags          []string
	Level         int
	NumOfHeadings int
	IsBookmark    bool
	IsTodo        bool
	IsDump        bool
	EditFile      bool
	Global        bool
	Quiet         bool
}

func (c Config) NoteType() string {
	if c.IsBookmark {
		return note.Bookmark
	} else if c.IsDump {
		return note.Dump
	} else if c.IsTodo {
		return note.Todo
	}
	return "dump"
}

func (c Config) Equals(other Config) bool {
	return c.Content == other.Content &&
		c.Description == other.Description &&
		c.Notespath == other.Notespath &&
		c.NoteType() == other.NoteType() &&
		slices.Equal(c.Tags, other.Tags) &&
		c.Level == other.Level &&
		c.NumOfHeadings == other.NumOfHeadings &&
		c.EditFile == other.EditFile &&
		c.Global == other.Global &&
		c.Quiet == other.Quiet
}
