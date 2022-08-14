package asyncjob

import (
	"context"
	"errors"
	"github.com/200Lab-Education/go-sdk/logger"
	"time"
)

var ErrTaskFailed = errors.New("job has failed after many retries")

var defaultRetryDurations = []time.Duration{
	time.Second * 5,
	time.Second * 15,
	time.Minute,
	time.Minute * 3,
}

type State int

const (
	Init State = iota
	Doing
	Retrying
	Failed
	Completed
)

type asyncJob struct {
	name           string
	logger         logger.Logger
	state          State
	handler        func(ctx context.Context) error
	retryDurations []time.Duration
	retryIndex     int
	cancelFunc     func()
	doneChan       chan interface{}
}

func NewAsyncJob(name string, logger logger.Logger, handler func(ctx context.Context) error) *asyncJob {
	return &asyncJob{
		name:           name,
		logger:         logger,
		handler:        handler,
		state:          Init,
		retryDurations: defaultRetryDurations,
		retryIndex:     -1,
	}
}

func (as *asyncJob) Execute(ctx context.Context) error {
	if as.state != Retrying {
		as.log("doing job: " + as.name)
		as.state = Doing
	}

	ctx, cancelF := context.WithCancel(ctx)
	as.cancelFunc = cancelF

	if err := as.handler(ctx); err != nil {
		return err
	}

	as.log("completed job: " + as.name)
	as.state = Completed
	return nil
}

func (as *asyncJob) Retry(ctx context.Context) error {
	if as.retryIndex == len(as.retryDurations)-1 {
		as.state = Failed
		as.log("failed task " + as.name)
		return ErrTaskFailed
	}

	as.state = Retrying
	as.retryIndex += 1

	as.log("prepare to retry job: "+as.name+" after", as.retryDurations[as.retryIndex])
	time.Sleep(as.retryDurations[as.retryIndex])
	as.log("retrying job: " + as.name)

	return as.Execute(ctx)
}

func (as *asyncJob) Cancel() {
	if as.cancelFunc != nil {
		as.cancelFunc()
	}
}

func (as *asyncJob) SetRetryDurations(retryDurations []time.Duration) {
	as.retryDurations = retryDurations
}

func (as *asyncJob) log(v ...interface{}) {
	if as.logger == nil {
		return
	}
	as.logger.Debugln(v)
}
