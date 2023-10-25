package app

import (
	"context"
	"sync"
	"time"

	"bharvest.io/axelmon/log"
)

type Monfunc func(ctx context.Context) (error)

func Run(ctx context.Context, c *Config) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)

	monitoringFuncs := []Monfunc{
		c.checkMaintainers,
		c.checkHeartbeats,
		c.checkEVMVotes,
	}

	wg := sync.WaitGroup{}
	wg.Add(len(monitoringFuncs))

	for _, f := range monitoringFuncs {
		go func(monFunc Monfunc) {
			defer wg.Done()

			err := monFunc(ctx)
			if err != nil {
				log.Error(err)
				return
			}
		}(f)
	}

	wg.Wait()
	cancel()

	return
}
