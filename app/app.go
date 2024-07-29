package app

import (
	"context"
	"sync"
	"time"

	"bharvest.io/axelmon/log"
)

type Monfunc func(ctx context.Context) error

func Run(ctx context.Context, c *Config) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	var monitoringFuncs []Monfunc

	if len(c.General.TargetSvcs) == 0 {
		monitoringFuncs = []Monfunc{c.checkMaintainers, c.checkHeartbeats, c.checkEVMVotes}
	} else {
		for _, targetSvc := range c.General.TargetSvcs {
			switch targetSvc {
			case MaintainerTargetSvc:
				monitoringFuncs = append(monitoringFuncs, c.checkMaintainers)
				break
			case HeartbeatTargetSvc:
				monitoringFuncs = append(monitoringFuncs, c.checkHeartbeats)
				break
			case EVMVoteTargetSvc:
				monitoringFuncs = append(monitoringFuncs, c.checkEVMVotes)
				break
			case VMVoteTargetSvc:
				monitoringFuncs = append(monitoringFuncs, c.checkVMVotes)
				break
			}
		}
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

	return
}
