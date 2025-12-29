package main

import (
	"fmt"
	"os"

	Config "github.com/jcfullmer/gatoRSS/internal/config"
)

func main() {
	configState, _ := Config.GetState()
	cmds := &Config.Commands{
		Handlers: map[string]func(*Config.State, Config.Command) error{
			"login": Config.HandlerLogin,
		},
	}
	if len(os.Args) < 2 {
		fmt.Print("No args")
		os.Exit(1)
	}
	argInput := os.Args

	cmd := Config.CreateCommand(argInput)
	fmt.Println(cmd.Name)
	err := cmds.Run(configState, cmd)
	if err != nil {
		fmt.Print(err)
	}
}
