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

func HandlerAddFeed(s *State, cmd Command, usr database.User) error {
	user, err := s.Db.GetUser(context.Background(), s.Cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("error getting data from database. Error: %v", err)
	}

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	id := uuid.New()
	name := cmd.Args[0]
	created_at := time.Now().UTC()
	url := cmd.Args[1]
	user_id := user.ID

	_, err = s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        id,
		Name:      name,
		CreatedAt: created_at,
		UpdatedAt: created_at,
		Url:       url,
		UserID:    user_id,
	})

	if err != nil {
		return fmt.Errorf("error with creation of new feed. Error: %v", err)
	}

	_, err = followFeed(s, url, usr)
	if err != nil {
		return err
	}

	fmt.Println("successfully created feed")
	fmt.Println(name)

	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting data from database. Error: %v", err)
	}
	for _, feed := range feeds {
		user, err := s.Db.GetUserFromID(context.Background(), feed.UserID)

		if err != nil {
			return fmt.Errorf("error getting data from database. Error: %v", err)
		}
		fmt.Println(len(feeds), "feeds in total")
		fmt.Println("Feed Name:")
		fmt.Println(feed.Name)
		fmt.Println("URL:")
		fmt.Println(feed.Url)
		fmt.Println("Created By:")
		fmt.Println(user.Name)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command, usr database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	selectedFeed, err := followFeed(s, cmd.Args[0], usr)
	if err != nil {
		return err
	}

	fmt.Println("Followed Feed:")
	fmt.Println(selectedFeed.Name)
	fmt.Println("Current User:")
	fmt.Println(usr.Name)

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	user, err := s.Db.GetUser(context.Background(), s.Cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("error getting data from database. Error: %v", err)
	}

	feedfollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting data from database. Error: %v", err)
	}

	if len(feedfollows) == 0 {
		fmt.Printf("Current User %v doesn't follow any feeds\n", user.Name)
		return nil
	}

	fmt.Printf("User: %v follows the following feeds;\n", user.Name)
	for _, feedfollow := range feedfollows {
		fmt.Println(feedfollow.FeedName)
	}

	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	feed, err := s.Db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	err = s.Db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
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

func followFeed(s *State, url string, usr database.User) (database.Feed, error) {
	selectedFeed, err := s.Db.GetFeedByURL(context.Background(), url)
	if err != nil {
		var feed database.Feed
		return feed, fmt.Errorf("error getting data from database. Error: %v", err)
	}

	id := uuid.New()
	currentTime := time.Now().UTC()

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    usr.ID,
		FeedID:    selectedFeed.ID,
	})

	if err != nil {
		var feed database.Feed
		return feed, fmt.Errorf("error with creation of new feed follow. Error: %v", err)
	}
	return selectedFeed, nil
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.Current_user_name)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
