package notes

import "errors"

type Config struct {
	Filepath      string
	Level         int
	NumOfHeadings int
	Help          bool
	OpenEditor    bool
}

func NewConfig(
	filepath string,
	level, numOfHeadings int,
	help, openEditor bool,
) (*Config, error) {
	if level < 1 || level > 6 {
		return nil, errors.New("")
	}
	c := new(Config)
	c.Filepath = filepath
	c.OpenEditor = openEditor
	c.Help = help
	c.NumOfHeadings = numOfHeadings
	c.Level = level
	return c, nil
}
