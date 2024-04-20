package note

import (
	"errors"
)

type Config struct {
	Content   string
	Mode      string
	Notespath string
	EditFile  bool
	Quiet     bool
}

func (c *Config) Validate() error {
	if c.Content == "" && !c.EditFile {
		return errors.New("nothing to note here")
	}
	return nil
}
