package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"log"
	"fmt"
	"errors"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type state struct {
	cfg *Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	index map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if c.index[cmd.name] == nil {
		return errors.New("Unknown command.")
	} 
	err := c.index[cmd.name]
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.index[name] != nil {
		fmt.Println("Command already exists.")
	}
	c.index[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if cmd.args == nil {
		return errors.New("No commands given.")
	}
	err := s.cfg.SetUser(s.cfg.CurrentUserName)
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}
	fmt.Printf("User has been. Welcome: %+v\n", s.cfg.CurrentUserName)
	return nil
}


func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(home, configFileName)
	return fullPath, nil
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
