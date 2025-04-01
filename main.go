package main

import (
	"database/sql"
	"fmt"
	"os"

	cnfg "github.com/MrR0b0t1001/aggregator/internal/config"
	dbpk "github.com/MrR0b0t1001/aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	config := cnfg.Read()

	db, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		fmt.Println("Database connection failed")
		os.Exit(1)
	}

	state := &cnfg.State{
		CurrState: &config,
		DB:        dbpk.New(db),
	}

	c := cnfg.Commands{
		CommandsMap: make(map[string]func(*cnfg.State, cnfg.Command) error),
	}

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments were provided")
		os.Exit(1)
	} else if len(os.Args) < 3 {
		fmt.Println("Username is required")
		os.Exit(1)
	}

	c.Register("login", cnfg.HandlerLogin)
	c.Register("register", cnfg.HandlerRegister)

	cmd := cnfg.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	c.Run(state, cmd)
}
