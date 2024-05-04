// Package views provides forms and other input TUI for notes app
package views

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/chaitanyabsprip/note/cmd/note/config"
)

// GetIssueConfiguration function  
func GetIssueConfiguration() (*config.Config, error) {
	c := &config.Config{IsIssue: true}
	tags := ""
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Placeholder("What is it about?").
				Prompt("▍").
				Value(&c.Title),
			huh.NewText().
				Title("Description").
				Placeholder("Describe your issue").
				Value(&c.Content).
				WithHeight(4),
			huh.NewInput().
				Title("Tags").
				Suggestions([]string{"bug", "task", "enhancement"}).
				Prompt("▍").
				Value(&tags),
		),
	).WithTheme(ThemeRosepine()).Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			os.Exit(130)
		}
		fmt.Println(err)
		os.Exit(1)
	}
	c.Tags = strings.Split(tags, ",")
	return c, nil
}
