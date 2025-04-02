package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	cmd "github.com/MrR0b0t1001/aggregator/internal/commands"
	cnfg "github.com/MrR0b0t1001/aggregator/internal/config"
	dbpk "github.com/MrR0b0t1001/aggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	config := cnfg.Read()

	state, err := cnfg.NewState()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		fmt.Println("Database connection failed")
		os.Exit(1)
	}

	state.DB = dbpk.New(db)

	commands := cmd.NewCommands()

	command := cmd.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := commands.Run(state, command); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
