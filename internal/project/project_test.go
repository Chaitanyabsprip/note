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

	projectName := "Test Project"
	projectPath := "/path/to/project"
	projectURL := "http://example.com"

	project, err := pm.AddProject(
		projectName,
		projectPath,
		projectURL,
	)
	if err != nil {
		t.Fatalf("Error adding project: %v", err)
	}

	if project.Name != projectName {
		t.Errorf("Expected project name %s, got %s", projectName, project.Name)
	}
	if project.Path != projectPath {
		t.Errorf("Expected project path %s, got %s", projectPath, project.Path)
	}
	if project.URL != projectURL {
		t.Errorf("Expected project URL %s, got %s", projectURL, project.URL)
	}

	// Try to add a project with the same name
	_, err = pm.AddProject(projectName, projectPath, projectURL)
	if err == nil {
		t.Error("Expected error when adding a project with the same name, got nil")
	}
	if !AlreadyExists(err) {
		t.Errorf("Expected AlreadyExists error, got %v", err)
	}
}

func TestUpdateProject(t *testing.T) {
	filepath := "test_projects.json"
	defer os.Remove(filepath)

	pm, err := NewProjectRepository(filepath)
	if err != nil {
		t.Fatalf("Error creating project manager: %v", err)
	}

	project, err := pm.AddProject(
		"Test Project",
		"/path/to/project",
		"http://example.com",
	)
	if err != nil {
		t.Fatalf("Error adding project: %v", err)
	}

	updatedProject, err := pm.UpdateProject(
		project.ID,
		"Updated Project",
		"/new/path",
		"http://updated.com",
	)
	if err != nil {
		t.Fatalf("Error updating project: %v", err)
	}

	if updatedProject.Name != "Updated Project" {
		t.Errorf(
			"Expected project name to be 'Updated Project', got '%s'",
			updatedProject.Name,
		)
	}
	if updatedProject.Path != "/new/path" {
		t.Errorf(
			"Expected project path to be '/new/path', got '%s'",
			updatedProject.Path,
		)
	}
	if updatedProject.URL != "http://updated.com" {
		t.Errorf(
			"Expected project URL to be 'http://updated.com', got '%s'",
			updatedProject.URL,
		)
	}
}
