// Package config provides config
package config

import (
	"slices"

	"github.com/chaitanyabsprip/note/internal/note"
)

// Config struct  
type Config struct {
	NoteType      string
	Content       string
	Description   string
	Notespath     string
	Project       string
	Title         string
	Status        note.Status
	Tags          []string
	NumOfHeadings int
	Level         int
	EditFile      bool
	Peek          bool
	Quiet         bool
}

// Equals method  
func (c Config) Equals(other Config) bool {
	return c.Content == other.Content &&
		c.Peek == other.Peek &&
		c.NoteType == other.NoteType &&
		c.Description == other.Description &&
		c.EditFile == other.EditFile &&
		c.Level == other.Level &&
		c.Notespath == other.Notespath &&
		c.NumOfHeadings == other.NumOfHeadings &&
		slices.Equal(c.Tags, other.Tags) &&
		c.Title == other.Title &&
		c.Status == other.Status &&
		c.Quiet == other.Quiet
}
