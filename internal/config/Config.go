package Config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

const configFileName = ".gatorconfig.json"

func get_config_path() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home_dir + "/" + configFileName, nil
}

type Config struct {
	DbURL           string    `json:"db_url"`
	CurrentUserName string    `json:"current_user_name"`
	CurrentUserID   uuid.UUID `json:"current_user_id"`
}

func Read() (Config, error) {
	conf_path, err := get_config_path()
	conf := Config{}
	if err != nil {
		return conf, err
	}
	conf_bytes, err := os.ReadFile(conf_path)
	if err != nil {
		return conf, err
	}
	json.Unmarshal(conf_bytes, &conf)
	return conf, nil
}

func writeConfig(c Config) error {
	conf_path, err := get_config_path()
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(conf_path, jsonData, os.FileMode(0777))
	if err != nil {
		return err
	}
	return nil
}

func (c Config) SetUser(username string, id uuid.UUID) error {
	c.CurrentUserName = username
	c.CurrentUserID = id
	err := writeConfig(c)
	if err != nil {
		return err
	}
	fmt.Printf("Set user as %v With an ID of %v\n", c.CurrentUserName, c.CurrentUserID)
	return nil
}
