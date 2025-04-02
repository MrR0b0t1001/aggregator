package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	cnfg "github.com/MrR0b0t1001/aggregator/internal/config"

	dbpk "github.com/MrR0b0t1001/aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *cnfg.State, cmd Command) error {
	user, err := s.DB.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User does not exist")
			os.Exit(1)
		}

		return err
	}

	if err := s.CurrState.SetUser(user.Name); err != nil {
		return err
	}

	log.Print("User has been set")

	return nil
}

func HandlerRegister(s *cnfg.State, cmd Command) error {
	user, err := s.DB.CreateUser(context.Background(), dbpk.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	})
	if err != nil {
		fmt.Println("User Not created")
		os.Exit(1)
	}

	fmt.Println("User Created")

	s.CurrState.SetUser(user.Name)
	log.Println(user)

	return nil
}

func HandlerReset(s *cnfg.State, cmd Command) error {
	if err := s.DB.DeleteAllUsers(context.Background()); err != nil {
		log.Println("Users reset failed")
		return err
	}

	log.Println("Users reset successful")
	return nil
}

func HandlerUsers(s *cnfg.State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if s.CurrState.CurrentUserName == user.Name {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %v\n", user.Name)
	}

	return nil
}
