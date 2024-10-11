package app

import (
	"bharvest.io/axelmon/server"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"bharvest.io/axelmon/log"
)

type Monfunc func(ctx context.Context) error

func Run(ctx context.Context, c *Config) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	c.alertChan = make(chan *alertMsg)

	go func() {
		for {
			select {
			case alert := <-c.alertChan:
				go func(msg *alertMsg) {
					var e error
					e = notifyTg(msg)
					if e != nil {
						log.Error(errors.New(fmt.Sprintf("error sending alert to telegram %v", e)))
					}
					e = notifySlack(msg)
					if e != nil {
						log.Error(errors.New(fmt.Sprintf("error sending alert to slack %v", e)))
					}
				}(alert)
			case <-ctx.Done():
				return
			}
		}
	}()

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

func SaveOnExit(stateFile string) {

	saveState := func() {
		log.Info("saving state...")
		//#nosec -- variable specified on command line
		f, e := os.OpenFile(stateFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if e != nil {
			log.Error(e)
			return
		}

		b, e := json.Marshal(&server.GlobalState)
		if e != nil {
			log.Error(e)
			return
		}
		_, _ = f.Write(b)
		_ = f.Close()
		log.Info("Axelmon exiting.")
	}
	saveState()
	//for {
	//	select {
	//	case <-ctx.Done():
	//		saveState()
	//		return
	//	case <-quitting:
	//		saveState()
	//		return
	//	}
	//}
}
