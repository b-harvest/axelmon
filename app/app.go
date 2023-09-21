package app

import (
	"context"

	"bharvest.io/axelmon/log"
)

func Run(ctx context.Context, c *Config) error {
	err := c.checkMaintainers(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	err = c.checkHeartbeat(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
