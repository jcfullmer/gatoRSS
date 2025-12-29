package Config

import (
	"fmt"
	"os"
)

type State struct {
	config *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func GetState() (*State, error) {
	config, err := Read()
	if err != nil {
		return &State{}, err
	}
	return &State{
		&config,
	}, nil
}

func CreateCommand(args []string) Command {
	name := args[1]
	var argies []string
	if len(args) >= 3 {
		argies = args[2:]
	}
	return Command{
		Name: name,
		Args: argies,
	}
}

func (c *Commands) Run(s *State, cmd Command) error {
	if _, ok := c.Handlers[cmd.Name]; !ok {
		return fmt.Errorf("command not found")
	}
	err := c.Handlers[cmd.Name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) error {
	if _, ok := c.Handlers[name]; ok {
		return fmt.Errorf("command already exists")
	}
	c.Handlers[name] = f
	return nil
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Print("a username is required")
		os.Exit(1)
		return fmt.Errorf("Please input a username to login.")
	}
	newName := cmd.Args[0]
	err := s.config.SetUser(newName)
	if err != nil {
		return err
	}
	fmt.Print("The user has been set.")
	return nil
}
