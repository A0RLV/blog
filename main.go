package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"errors"
	
	"github.com/A0RLV/blog/internal/config"
)


type state struct {
	cfg Config
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

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	s := state{cfg}
	cs := commands{}

	cs.register("login", handlerLogin(*state, command))
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: <command> <arguments>")
		os.Exit(1)
	}
	cmd_input := os.Args[1]

	cmd_name := strings.Split(cmd_input, " ")[0]
	cmd_args := strings.Split(cmd_input, " ")[1:]
	c := command{cmd_name, cmd_args}

	err := cs.run(s, c)
	if err != nil {
		fmt.Fatalf("%v", err)
	}

	// err = cfg.SetUser("kaz")
	// if err != nil {
	// 	log.Fatalf("couldn't set current user: %v", err)
	// }

	// cfg, err = config.Read()
	// if err != nil {
	// 	log.Fatalf("error reading config: %v", err)
	// }
	// fmt.Printf("Read config again: %+v\n", cfg)


}
