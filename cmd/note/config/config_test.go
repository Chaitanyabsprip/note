// Package config provides config
package config

import (
	"testing"

	"github.com/chaitanyabsprip/note/internal/note"
)

func TestConfig_NoteType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		config   Config
	}{
		{
			name:     "Bookmark NoteType",
			expected: note.Bookmark,
			config:   Config{NoteType: note.Bookmark},
		},
		{
			name:     "Dump NoteType",
			expected: note.Dump,
			config:   Config{NoteType: note.Dump},
		},
		{
			name:     "Todo NoteType",
			expected: note.Todo,
			config:   Config{NoteType: note.Todo},
		},
		{
			name:     "Issue NoteType",
			expected: note.Issue,
			config:   Config{NoteType: note.Issue},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.config
			if got := c.NoteType; got != tt.expected {
				t.Errorf("Config.NoteType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConfig_Equals(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		other    Config
		expected bool
	}{
		{
			name: "Equal Configs",
			config: Config{
				Content:       "Test content",
				Description:   "Test description",
				Notespath:     "/path/to/notes",
				Project:       "Test project",
				Title:         "Test title",
				Tags:          []string{"tag1", "tag2"},
				Level:         2,
				NumOfHeadings: 3,
				EditFile:      true,
				NoteType:      note.Bookmark,
				Quiet:         false,
			},
			other: Config{
				Content:       "Test content",
				Description:   "Test description",
				Notespath:     "/path/to/notes",
				Project:       "Test project",
				Title:         "Test title",
				Tags:          []string{"tag1", "tag2"},
				Level:         2,
				NumOfHeadings: 3,
				EditFile:      true,
				NoteType:      note.Bookmark,
				Quiet:         false,
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.config
			if got := c.Equals(tt.other); got != tt.expected {
				t.Errorf("Config.Equals() = %v, want %v", got, tt.expected)
			}
		})
	}
}
