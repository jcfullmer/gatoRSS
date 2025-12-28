package Config

import (
	"errors"
)

type State struct {
	config *Config
}

type command struct {
	name string
	args []string
}

```func handlerLogin(s *State, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Please input an argument for the command.")
	}
}
```
