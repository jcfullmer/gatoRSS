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
		Handlers: map[string]func(*Config.State, Config.Command) error{},
	}
	cmds.Register("login", Config.HandlerLogin)
	cmds.Register("register", Config.HandlerRegister)
	cmds.Register("reset", Config.HandlerReset)
	cmds.Register("users", Config.HandlerUsers)
	cmds.Register("agg", Config.HandlerAgg)
	cmds.Register("feeds", Config.HandlerFeeds)
	cmds.Register("addfeed", Config.MiddlewareLoggedIn(Config.HandlerAddFeed))
	cmds.Register("follow", Config.MiddlewareLoggedIn(Config.HandlerFollow))
	cmds.Register("following", Config.MiddlewareLoggedIn(Config.HandlerFollowing))
	cmds.Register("unfollow", Config.MiddlewareLoggedIn(Config.HandlerUnfollow))
	cmds.Register("browse", Config.HandlerBrowse)
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
