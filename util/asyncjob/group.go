package asyncjob

import (
	"context"
	"errors"
	"fmt"
	"github.com/200Lab-Education/go-sdk/logger"
	"sync"
	"time"
)

type Job interface {
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	Cancel()
	SetRetryDurations(retryDurations []time.Duration)
}

type asyncGroup struct {
	jobs       []Job
	wg         *sync.WaitGroup
	isParallel bool
	logger     logger.Logger
}

func Compose(isParallel bool, logger logger.Logger, jobs ...Job) *asyncGroup {
	return &asyncGroup{
		jobs:       jobs,
		isParallel: isParallel,
		logger:     logger,
		wg:         new(sync.WaitGroup),
	}
}

func (ag *asyncGroup) Run(ctx context.Context) error {
	ag.wg.Add(len(ag.jobs))
	var err error

	for _, t := range ag.jobs {
		if ag.isParallel {
			go func(as Job) {
				defer func() {
					if err := recover(); err != nil {
						j := t.(*asyncJob)
						ag.logger.Error(j.name, err)

						err = errors.New(fmt.Sprintf("%v", err))
						ag.wg.Done()
					}
				}()

				err = ag.do(ctx, as)
				ag.wg.Done()
			}(t)
		} else {
			if err := ag.do(ctx, t); err != nil {
				return err
			}
		}
	}

	if ag.isParallel {
		ag.wg.Wait()
	}
	return err
}

func (ag *asyncGroup) do(ctx context.Context, as Job) error {
	if err := as.Execute(ctx); err != nil {
		for {
			if err := as.Retry(ctx); err != nil {
				if err == ErrTaskFailed {
					return err
				}
				continue
			}
			return nil
		}
	}
	return nil
}
