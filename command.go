package main

import "fmt"

type command struct {
	Name string
	Args []string
}
type commands struct {
	Commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	cmdFunc, exists := c.Commands[cmd.Name]
	if !exists {
		return fmt.Errorf("cmd with name:%s has not been registerd", cmd.Name)
	}
	err := cmdFunc(s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.Commands[name] = f
}
