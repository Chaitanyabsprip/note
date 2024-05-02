package note

import (
	"strings"
	"testing"
	"time"
)

func TestNoteTypeToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		noteType noteType
		content  string
		expected string
	}{
		// Normal Cases
		{
			name:     "BookmarkCreation",
			noteType: &bookmark{description: "Bookmark description", tags: []string{"tag1", "tag2"}},
			content:  "https://example.com",
			expected: "[Example Domain](https://example.com)  \ntags: tag1, tag2  \nBookmark description",
		},
		{
			name:     "NotesCreation",
			noteType: new(notes),
			content:  "This is a test note.",
			expected: "This is a test note.",
		},
		{
			name:     "IssueCreation",
			noteType: NewIssue("Test Issue", "This is a test issue description.", []string{"bug", "enhancement"}, time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)),
			content:  "This is a test issue.",
			expected: "## Test Issue\n\ncreatedAt: Sat Jan  1 00:00:00 UTC 2022\nstatus: Open\nlabels: bug, enhancement\n\nThis is a test issue.\n\n### Comments\n---",
		},
		{
			name:     "TodoCreation",
			noteType: new(todo),
			content:  "This is a test todo.",
			expected: "- [ ] This is a test todo.",
		},
		// Edge Cases
		{
			name:     "EmptyContent",
			noteType: new(bookmark),
			content:  "",
			expected: "[]()  \ntags:",
		},
		{
			name:     "InvalidURL",
			noteType: new(bookmark),
			content:  "invalid-url",
			expected: "[](invalid-url)  \ntags:\n",
		},
		{
			name:     "EmptyLabels",
			noteType: NewIssue("", "", []string{}, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)),
			content:  "This is a test issue.",
			expected: "## \n\ncreatedAt: \nstatus: Open\nlabels: \n\nThis is a test issue.\n\n### Comments\n---",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			md, err := tc.noteType.toMarkdown(tc.content)
			if err != nil {
				t.Errorf("Error converting %s to Markdown: %v", tc.name, err)
			}
			if strings.TrimSpace(md) != strings.TrimSpace(tc.expected) {
				t.Errorf("%s Markdown does not match expected. Got: %s, Expected: %s", tc.name, md, tc.expected)
			}
		})
	}
}
