package main

import (
	"context"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no arguments provided.usage:%s <user name>", cmd.Name)
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments.usage: %s <user name>", cmd.Name)
	}
	userName := cmd.Args[0]
	dbUser, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("Error getting username:%s from DB. ERR:%v", userName, err)
	}
	err = s.Conf.SetUser(dbUser.Name)
	if err != nil {
		return err
	}
	fmt.Printf("the user name has been set to %s", userName)
	return nil
}
