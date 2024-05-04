// Package config provides config
package config

import (
	"testing"

	"github.com/chaitanyabsprip/note/pkg/note"
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
			config:   Config{IsBookmark: true},
		},
		{
			name:     "Dump NoteType",
			expected: note.Dump,
			config:   Config{IsDump: true},
		},
		{
			name:     "Todo NoteType",
			expected: note.Todo,
			config:   Config{IsTodo: true},
		},
		{
			name:     "Issue NoteType",
			expected: note.Issue,
			config:   Config{IsIssue: true},
		},
		{
			name:     "Default NoteType",
			expected: "dump",
			config:   Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.config
			if got := c.NoteType(); got != tt.expected {
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
				IsBookmark:    true,
				IsDump:        false,
				IsIssue:       false,
				IsTodo:        false,
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
				IsBookmark:    true,
				IsDump:        false,
				IsIssue:       false,
				IsTodo:        false,
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
