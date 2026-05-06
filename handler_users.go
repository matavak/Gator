package main

import (
	"context"
	"fmt"
)

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.Args) >= 1 {
		return fmt.Errorf("too many arguments.usage: %s ", cmd.Name)
	}
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return err
	}
	for _, v := range users {
		userName := v.Name
		if s.Conf.CurrentUserName == userName {
			userName += " (current)"
		}
		fmt.Printf("* %s \n", userName)
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no arguments provided.usage:%s <user name>", cmd.Name)
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments.usage: %s <user name>", cmd.Name)
	}
	userName := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), userName)
	if err != nil {
		return err
	}
	err = s.Conf.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("the user has been registered to DB  %v \n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("too many arguments.usage: %s ", cmd.Name)
	}
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Could not reset user table in db ERR:%v", err)
	}
	fmt.Printf("user table has been reset")
	return nil
}
