package commands

import (
	"errors"
	// Your imports
	"github.com/MrR0b0t1001/aggregator/internal/config"
	// other imports
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandsMap map[string]func(*config.State, Command) error
}

func NewCommands() *Commands {
	cmds := &Commands{
		CommandsMap: make(map[string]func(*config.State, Command) error),
	}

	// Register all commands
	cmds.Register("login", HandlerLogin)
	cmds.Register("register", HandlerRegister)
	cmds.Register("reset", HandlerReset)
	cmds.Register("users", HandlerUsers)
	cmds.Register("agg", HandlerAgg)

	return cmds
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
	c.CommandsMap[name] = f
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	f, ok := c.CommandsMap[cmd.Name]
	if !ok {
		return errors.New("Command does not exist")
	}

	f(s, cmd)

	return nil
}
