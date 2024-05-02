package main

import (
	"slices"

	"github.com/chaitanyabsprip/note/pkg/note"
)

type Config struct {
	Status        note.Status
	Content       string
	Description   string
	Notespath     string
	Title         string
	Tags          []string
	Level         int
	NumOfHeadings int
	EditFile      bool
	Global        bool
	IsBookmark    bool
	IsDump        bool
	IsIssue       bool
	IsTodo        bool
	Quiet         bool
}

func (c Config) NoteType() string {
	if c.IsBookmark {
		return note.Bookmark
	} else if c.IsDump {
		return note.Dump
	} else if c.IsTodo {
		return note.Todo
	} else if c.IsIssue {
		return note.Issue
	}
	return "dump"
}

func (c Config) Equals(other Config) bool {
	return c.Content == other.Content &&
		c.Description == other.Description &&
		c.EditFile == other.EditFile &&
		c.Global == other.Global &&
		c.Level == other.Level &&
		c.NoteType() == other.NoteType() &&
		c.Notespath == other.Notespath &&
		c.NumOfHeadings == other.NumOfHeadings &&
		slices.Equal(c.Tags, other.Tags) &&
		c.Title == other.Title &&
		c.Status == other.Status &&
		c.Quiet == other.Quiet
}
