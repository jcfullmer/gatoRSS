package main

import (
	"fmt"

	Config "github.com/jcfullmer/gatoRSS/internal/config"
)

func main() {
	config, _ := Config.Read()
	config.SetUser("Claye")
	config, _ = Config.Read()
	fmt.Println(config)
}
