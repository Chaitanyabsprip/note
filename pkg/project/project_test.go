package project

import (
	"os"
	"testing"
)

func TestAddProject(t *testing.T) {
	filepath := "test_projects.json"
	defer os.Remove(filepath)

	pm, err := NewProjectRepository(filepath)
	if err != nil {
		t.Fatalf("Error creating project manager: %v", err)
	}

	_, err = pm.AddProject("Test Project", "/path/to/project", "http://example.com")
	if err != nil {
		t.Fatalf("Error adding project: %v", err)
	}

	if len(pm.projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(pm.projects))
	}
}

func TestUpdateProject(t *testing.T) {
	filepath := "test_projects.json"
	defer os.Remove(filepath)

	pm, err := NewProjectRepository(filepath)
	if err != nil {
		t.Fatalf("Error creating project manager: %v", err)
	}

	project, err := pm.AddProject("Test Project", "/path/to/project", "http://example.com")
	if err != nil {
		t.Fatalf("Error adding project: %v", err)
	}

	updatedProject, err := pm.UpdateProject(project.ID, "Updated Project", "/new/path", "http://updated.com")
	if err != nil {
		t.Fatalf("Error updating project: %v", err)
	}

	if updatedProject.Name != "Updated Project" {
		t.Errorf("Expected project name to be 'Updated Project', got '%s'", updatedProject.Name)
	}
	if updatedProject.Path != "/new/path" {
		t.Errorf("Expected project path to be '/new/path', got '%s'", updatedProject.Path)
	}
	if updatedProject.URL != "http://updated.com" {
		t.Errorf("Expected project URL to be 'http://updated.com', got '%s'", updatedProject.URL)
	}
}
