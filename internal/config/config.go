package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	dbpk "github.com/MrR0b0t1001/aggregator/internal/database"
	"github.com/google/uuid"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	CurrState *Config
	DB        *dbpk.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandsMap map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	_, ok := c.CommandsMap[name]
	if !ok {
		c.CommandsMap[name] = f
	}

	return
}

func (c *Commands) Run(s *State, cmd Command) error {
	f, ok := c.CommandsMap[cmd.Name]
	if !ok {
		return errors.New("Command does not exist")
	}

	f(s, cmd)

	return nil
}

/////////////////////////////////////////////////////////////////

func (c Config) SetUser(username string) error {
	c.CurrentUserName = username

	err := write(c)
	if err != nil {
		return err
	}

	return nil
}

func Read() Config {
	filepath, err := getConfigFilePath()
	if err != nil {
		return Config{}
	}

	file, err := os.Open(filepath)
	if err != nil {
		return Config{}
	}

	defer file.Close()

	var config Config

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}
	}

	return config
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filepath := filepath.Join(homeDir, configFileName)

	return filepath, nil
}

func write(cfg Config) error {
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return nil
	}

	defer file.Close()

	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if _, err := file.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Username not provided.")
	}

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

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Username not provided.")
	}

	arg := dbpk.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      os.Args[2],
	}

	user, err := s.DB.CreateUser(context.Background(), arg)
	if err != nil {
		fmt.Println("User Not created")
		os.Exit(1)
	}

	fmt.Println("User Created")

	s.CurrState.SetUser(user.Name)
	log.Println(user)

	return nil
}
