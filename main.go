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
	cfg *config.Config
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
	err := c.index[cmd.name](s, cmd)
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
	fmt.Printf("Command registered: %v\n", name)
}

func handlerLogin(s *state, cmd command) error {
	fmt.Printf("Initializing login...\n")
	if cmd.args == nil {
		return errors.New("No commands given.")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		log.Fatalf("Couldn't set current user: %v", err)
	}
	fmt.Printf("User login complete. Welcome: %+v\n", cmd.args[0])
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	s := state{&cfg}
	index := make(map[string]func(*state, command) error)
	cs := commands{index}

	cs.register("login", handlerLogin)
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: <command> <arguments>")
		os.Exit(1)
	}
	cmd_name := os.Args[1]
	cmd_args := strings.Split(os.Args[2], " ")
	// fmt.Printf("cmd_args = %v", cmd_args)
	c := command{cmd_name, cmd_args}

	err = cs.run(&s, c)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
