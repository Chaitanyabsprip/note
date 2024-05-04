// Package note provides note  
package note

import (
	"testing"
)

func TestNote_Note(t *testing.T) {
	type fields struct {
		Status      Status
		Content     string
		Description string
		NotesPath   string
		Title       string
		Type        string
		Tags        []string
		EditFile    bool
		HidePreview bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Note{
				Status:      tt.fields.Status,
				Content:     tt.fields.Content,
				Description: tt.fields.Description,
				NotesPath:   tt.fields.NotesPath,
				Title:       tt.fields.Title,
				Type:        tt.fields.Type,
				Tags:        tt.fields.Tags,
				EditFile:    tt.fields.EditFile,
				HidePreview: tt.fields.HidePreview,
			}
			if err := n.Note(); (err != nil) != tt.wantErr {
				t.Errorf("Note.Note() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
