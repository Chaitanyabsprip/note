// Package project provides project  
package project

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Project struct  
type Project struct {
	Name string `json:"name"`
	Path string `json:"path"`
	URL  string `json:"url"`
	ID   int    `json:"id"`
}

// Repository interface  
type Repository interface {
	GetProject(name string) *Project
	AddProject(name, path, url string) (*Project, error)
	UpdateProject(id int, name, path, url string) (*Project, error)
}

type repositoryImpl struct {
	configPath string
	projects   []*Project
	mu         sync.Mutex
}

// NewProjectRepository function  
func NewProjectRepository(configPath string) (Repository, error) {
	pr := &repositoryImpl{configPath: configPath}
	err := pr.loadProjects()
	if err != nil {
		return nil, err
	}
	return pr, nil
}

// GetProject method  
func (pr *repositoryImpl) GetProject(name string) *Project {
	for _, project := range pr.projects {
		if project.Name == name {
			return project
		}
	}
	return nil
}

// AddProject method  
func (pr *repositoryImpl) AddProject(
	name, path, url string,
) (*Project, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	for _, p := range pr.projects {
		if p.Name == name {
			return nil, errors.New("project with same name already exists")
		}
	}
	id := 0
	if len(pr.projects) > 0 {
		id = pr.projects[len(pr.projects)-1].ID + 1
	}
	// if url == "" {
	// check whether path is a git directory, if so get remote url from that
	// }
	project := &Project{
		ID:   id,
		Name: name,
		Path: path,
		URL:  url,
	}
	pr.projects = append(pr.projects, project)
	err := pr.saveProjects()
	if err != nil {
		return nil, err
	}
	return project, nil
}

// UpdateProject method  
func (pr *repositoryImpl) UpdateProject(
	id int,
	name, path, url string,
) (*Project, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	for _, p := range pr.projects {
		if p.ID == id {
			p.Name = name
			p.Path = path
			p.URL = url
			err := pr.saveProjects()
			if err != nil {
				return nil, err
			}
			return p, nil
		}
	}
	return nil, errors.New("project not found")
}

// GetRepositoryRoot uses git CLI to find the root of the repository.
func GetRepositoryRoot(dirpath string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = dirpath
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(out.String(), "\n", "")
}

func (pr *repositoryImpl) loadProjects() error {
	if _, err := os.Stat(pr.configPath); os.IsNotExist(err) {
		err = os.WriteFile(pr.configPath, []byte("[]"), 0o644)
		if err != nil {
			return err
		}
		pr.projects = []*Project{}
		return nil
	}
	data, err := os.ReadFile(pr.configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pr.projects)
	if err != nil {
		return err
	}
	return nil
}

func (pr *repositoryImpl) saveProjects() error {
	data, err := json.MarshalIndent(pr.projects, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(pr.configPath, data, 0o644)
	if err != nil {
		return err
	}
	return nil
}

// AlreadyExists function  
func AlreadyExists(err error) bool {
	return err.Error() == "project with same name already exists"
}
