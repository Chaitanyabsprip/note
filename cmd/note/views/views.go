// Package views provides forms and other input TUI for notes app
package views

import (
	"errors"
	"fmt"
	"net/url"
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

// GetBookmarkConfiguration function  
func GetBookmarkConfiguration() (*config.Config, error) {
	c := &config.Config{IsBookmark: true}
	tags := ""
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Link").
				Prompt(": ").
				Validate(func(s string) error {
					if _, err := url.ParseRequestURI(s); err != nil {
						return errors.New("invalid URL")
					}
					return nil
				}).
				Inline(true).
				Value(&c.Content),
			huh.NewInput().
				Title("Tags").
				Prompt(": ").
				Value(&tags).Inline(true),
			huh.NewNote().Title(""),
			huh.NewText().
				Placeholder("Describe your bookmark").
				Value(&c.Description).
				WithHeight(2),
		),
	).WithHeight(7).WithTheme(ThemeRosepine()).Run()
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
