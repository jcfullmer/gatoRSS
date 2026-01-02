package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	Config "github.com/jcfullmer/gatoRSS/internal/config"

	database "github.com/jcfullmer/gatoRSS/internal/database"
)

func main() {
	configState, _ := Config.GetState()
	cmds := &Config.Commands{
		Handlers: map[string]func(*Config.State, Config.Command) error{
			"login":    Config.HandlerLogin,
			"register": Config.HandlerRegister,
			"reset":    Config.HandlerReset,
			"users":    Config.HandlerUsers,
			"agg":      Config.HandlerAgg,
		},
	}
	if len(os.Args) < 2 {
		fmt.Print("No args")
		os.Exit(1)
	}
	argInput := os.Args

	db, err := sql.Open("postgres", configState.Config.DbURL)
	dbQueries := database.New(db)
	configState.Db = dbQueries
	cmd := Config.CreateCommand(argInput)
	err = cmds.Run(configState, cmd)
	if err != nil {
		fmt.Print(err)
	}
}
