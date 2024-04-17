package notes

import (
	"errors"
)

type Config struct {
	Filepath      string
	Level         int
	NumOfHeadings int
	Help          bool
	OpenEditor    bool
}

func (c *Config) Validate() error {
	if c.Level > 6 || c.Level < 1 {
		return errors.New("level can be between [1, 6]")
	}
	return nil
}
