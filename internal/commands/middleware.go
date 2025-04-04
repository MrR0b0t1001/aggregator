package commands

import (
	"context"

	cnfg "github.com/MrR0b0t1001/aggregator/internal/config"
	"github.com/MrR0b0t1001/aggregator/internal/database"
)

func MiddlewareLoggedIn(
	handler func(s *cnfg.State, c Command, user database.User) error,
) func(s *cnfg.State, cmd Command) error {
	return func(s *cnfg.State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.CurrState.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
