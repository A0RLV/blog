package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/A0RLV/blog/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	s := state{cfg}
	cs := commands{}

	cs.register("login", config.handlerLogin())
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: <command> <arguments>")
		os.Exit(1)
	}
	cmd_input := os.Args[1]

	cmd_name := strings.Split(cmd_input, " ")[0]
	cmd_args := strings.Split(cmd_input, " ")[1:]
	c := command{cmd_name, cmd_args}

	err := c.config.run(s, c)
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
