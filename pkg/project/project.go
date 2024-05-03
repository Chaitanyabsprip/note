package project

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"sync"
)

type Project struct {
	Name string `json:"name"`
	Path string `json:"path"`
	URL  string `json:"url"`
	ID   int    `json:"id"`
}

type ProjectRepository struct {
	configPath string
	projects   []*Project
	mu         sync.Mutex
}

func NewProjectRepository(configPath string) (*ProjectRepository, error) {
	pr := &ProjectRepository{configPath: configPath}
	err := pr.loadProjects()
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (pr *ProjectRepository) GetProject(name string) *Project {
	for _, project := range pr.projects {
		if project.Name == name {
			return project
		}
	}
	return nil
}

func (pr *ProjectRepository) AddProject(name, path, url string) (*Project, error) {
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

func (pr *ProjectRepository) UpdateProject(id int, name, path, url string) (*Project, error) {
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

func GetRepositoryRoot(dirpath string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = dirpath
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return out.String()
}

func (pr *ProjectRepository) loadProjects() error {
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

func (pr *ProjectRepository) saveProjects() error {
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
