package cli

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/rss"
	"strconv"
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
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <duration>\n duration format = '1s', '1m', '1h'", cmd.Name)
	}
	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("incorrect duration format. duration format should be '1s', '1m', '1h'")
	}
	fmt.Printf("Collecting feeds every %v\n", duration)

	ticker := time.NewTicker(duration)

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func HandlerAddFeed(s *State, cmd Command, usr database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	id := uuid.New()
	name := cmd.Args[0]
	created_at := time.Now().UTC()
	url := cmd.Args[1]
	user_id := usr.ID

	_, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
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
	fmt.Println(len(feeds), "feeds in total")
	for _, feed := range feeds {
		user, err := s.Db.GetUserFromID(context.Background(), feed.UserID)

		if err != nil {
			return fmt.Errorf("error getting data from database. Error: %v", err)
		}
		fmt.Println("Feed Name:")
		fmt.Println(feed.Name)
		fmt.Println("URL:")
		fmt.Println(feed.Url)
		fmt.Println("Created By:")
		fmt.Println(user.Name)
		fmt.Println(feed.LastFetchedAt)
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

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	limit := 2
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <limit|optional>", cmd.Name)
	} else if len(cmd.Args) == 1 {
		tmpLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("error with command '%v'\n usage of command: %s <limit|optional>", err, cmd.Name)
		}
		limit = tmpLimit
	}

	posts, err := s.Db.GetPosts(context.Background(), database.GetPostsParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title:\n%v\n", post.Title)
		if post.PublishedAt.Valid == true {
			fmt.Printf("Published at:\n%v\n", post.PublishedAt.Time.Format(time.RFC1123))
		} else {
			fmt.Printf("Published at:\nUnknown\n")
		}
		fmt.Printf("URL:\n%v\n\n", post.Url)

		fmt.Printf("Content:\n%v\n", post.Description)

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

func scrapeFeeds(s *State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.Db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})

	if err != nil {
		return err
	}

	feedrss, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	feedrss.UnescapeStrings()
	for _, item := range feedrss.RSSChannel.Item {

		id := uuid.New()
		currentTime := time.Now().UTC()
		valid := true
		publishedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Println("could not parse pubDate", item.PubDate, "error", err)
			valid = false
		}

		_, err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          id,
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: sql.NullTime{
				Time:  publishedTime,
				Valid: valid,
			},
			FeedID: feed.ID,
		})
		if err != nil {
			return err
		}
		fmt.Printf("The following post successfully added\n%v\n", item.Title)

	}

	return nil
}
