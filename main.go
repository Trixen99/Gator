package main

import (
	"database/sql"
	"fmt"
	"gator/internal/cli"
	"gator/internal/config"
	"gator/internal/database"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	curState, commands, err := startup()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Println("You have not provided a command name you would like to run")
		os.Exit(1)
	}

	var command cli.Command

	if len(arguments) == 2 {
		command = cli.Command{
			Name: arguments[1],
		}
	} else {
		command = cli.Command{
			Name: arguments[1],
			Args: arguments[2:],
		}
	}

	err = commands.Run(&curState, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func startup() (cli.State, cli.Commands, error) {
	jsonconfig, err := config.Read()
	if err != nil {
		return cli.State{}, cli.Commands{}, fmt.Errorf("Problem reading config file: %w", err)
	}

	dbURL := jsonconfig.Db_url
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	curState := cli.State{
		Db:  dbQueries,
		Cfg: &jsonconfig,
	}

	commands := cli.Commands{
		Cmds: make(map[string]func(*cli.State, cli.Command) error),
	}

	commands.Register("login", cli.HandlerLogin)
	commands.Register("register", cli.HandlerRegister)
	commands.Register("reset", cli.HandlerReset)
	commands.Register("users", cli.HandlerUsers)
	commands.Register("agg", cli.HandlerAgg)
	commands.Register("addfeed", cli.MiddlewareLoggedIn(cli.HandlerAddFeed))
	commands.Register("feeds", cli.HandlerFeeds)
	commands.Register("follow", cli.MiddlewareLoggedIn(cli.HandlerFollow))
	commands.Register("following", cli.MiddlewareLoggedIn(cli.HandlerFollowing))
	commands.Register("unfollow", cli.MiddlewareLoggedIn(cli.HandlerUnfollow))

	return curState, commands, nil

}
