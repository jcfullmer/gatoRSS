package Config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	database "github.com/jcfullmer/gatoRSS/internal/database"
	"github.com/jcfullmer/gatoRSS/internal/rss"
)

type State struct {
	Db     *database.Queries
	Config *Config
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
		nil,
		&config,
	}, nil
}

func CreateCommand(args []string) Command {
	name := args[1]
	var argies []string
	if len(args) >= 2 {
		argies = args[2:]
	}
	return Command{
		Name: name,
		Args: argies,
	}
}

func (c *Commands) Run(s *State, cmd Command) error {
	if _, ok := c.Handlers[cmd.Name]; !ok {
		return fmt.Errorf("command %v not found", cmd.Name)
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
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println("User doesn't exist")
		os.Exit(1)
	}
	newName := cmd.Args[0]
	err = s.Config.SetUser(newName)
	if err != nil {
		return err
	}
	fmt.Print("The user has been set.")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Please input a name to register.")
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err == sql.ErrNoRows {
		dbUser := database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0],
		}

		user, err := s.Db.CreateUser(context.Background(), dbUser)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		Config.SetUser(*s.Config, dbUser.Name)
		fmt.Println("the user was created")
		fmt.Println(user)
		return nil
	} else if err != nil {
		fmt.Println("an error occured when getting user")
		os.Exit(1)
	} else {
		fmt.Println("user already exists")
		os.Exit(1)
	}
	return nil
}

func HandlerReset(s *State, _ Command) error {
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error when resetting databse: %v", err)
		os.Exit(1)
	}
	fmt.Println("Successfully reset database.")
	os.Exit(0)
	return nil
}

func HandlerUsers(s *State, _ Command) error {
	list, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error when Getting users: %v\n", err)
		os.Exit(1)
	}
	if len(list) == 0 {
		fmt.Println("No users in database.")
		os.Exit(1)
	}
	for _, user := range list {
		if user == s.Config.CurrentUserName {
			fmt.Printf("* %v (current)\n", user)
			continue
		}
		fmt.Printf("* %v\n", user)
	}
	os.Exit(0)
	return nil
}

func HandlerAgg(s *State, _ Command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	RSS, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		fmt.Println("Error with FetchFeed func")
		fmt.Println(err)
		return err
	}
	fmt.Println(RSS)
	return nil
}
