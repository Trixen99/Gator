package cli

import (
	"context"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/rss"
	"time"

	"github.com/google/uuid"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Cmds map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("name not found in database, error: %v", err)
	}

	err = s.Cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("User successfully updated")

	return nil
}

func HandlerRegister(s *State, cmd Command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	id := uuid.New()
	created_at := time.Now().UTC()
	updated_at := created_at
	name := cmd.Args[0]
	_, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name:      name,
	})

	if err != nil {
		return fmt.Errorf("error with creation of new user. Error: %v", err)
	}

	err = s.Cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("successfully created user")
	fmt.Println(name)

	return nil

}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error with resetting usuer table")
	}
	fmt.Println("All users deleted")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	allUsers, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting data from database. Error: %v", err)
	}

	for _, user := range allUsers {
		fmt.Print(user.Name)
		if user.Name == s.Cfg.Current_user_name {
			fmt.Print(" (current)")
		}
		fmt.Print("\n")

	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	data, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	fmt.Println(data)
	data.UnescapeStrings()

	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	_, ok := c.Cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("Command %v not found", cmd.Name)
	}

	err := c.Cmds[cmd.Name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Cmds[name] = f
}
