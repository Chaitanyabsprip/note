package main

import "slices"

type Config struct {
	Content       string
	Description   string
	Notespath     string
	Type          string
	Tags          []string
	Level         int
	NumOfHeadings int
	EditFile      bool
	Global        bool
	Quiet         bool
}

func (c Config) Equals(other Config) bool {
	return c.Content == other.Content &&
		c.Description == other.Description &&
		c.Notespath == other.Notespath &&
		c.Type == other.Type &&
		slices.Equal(c.Tags, other.Tags) &&
		c.Level == other.Level &&
		c.NumOfHeadings == other.NumOfHeadings &&
		c.EditFile == other.EditFile &&
		c.Global == other.Global &&
		c.Quiet == other.Quiet
}
