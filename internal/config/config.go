package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	dbpk "github.com/MrR0b0t1001/aggregator/internal/database"
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

func NewState() (*State, error) {
	state := &State{
		CurrState: &Config{},
		DB:        nil,
	}

	configPath := getConfigFilePath()

	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &state.CurrState); err != nil {
			return nil, err
		}
	}

	return state, nil
}

/////////////////////////////////////////////////////////////////

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

func Read() Config {
	filepath := getConfigFilePath()

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

func getConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Path incorrect")
		return ""
	}

	filepath := filepath.Join(homeDir, configFileName)

	return filepath
}

func write(config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getConfigFilePath(), data, 0)
}
