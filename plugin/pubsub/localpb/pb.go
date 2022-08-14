package localpb

import (
	"context"
	"flag"
	"fmt"
	"github.com/200Lab-Education/go-sdk/logger"
	pb "github.com/200Lab-Education/go-sdk/plugin/pubsub"
	"sync"
)

type pubsub struct {
	prefix       string
	locker       *sync.RWMutex
	logger       logger.Logger
	logEnabled   bool
	wg           *sync.WaitGroup
	gracefulStop bool
	messageQueue chan *pb.Event
	mapChannel   map[pb.Channel][]chan *pb.Event
	stopChan     chan bool
	isStopping   bool
}

func NewPubsub(prefix string) *pubsub {
	return &pubsub{
		locker:       new(sync.RWMutex),
		wg:           new(sync.WaitGroup),
		messageQueue: make(chan *pb.Event, 1000),
		mapChannel:   make(map[pb.Channel][]chan *pb.Event),
		stopChan:     make(chan bool),
		prefix:       prefix,
	}
}

func (ps *pubsub) GetPrefix() string {
	return ps.prefix
}

func (ps *pubsub) Get() interface{} {
	return ps
}

func (ps *pubsub) Name() string {
	return "pubsub"
}

func (ps *pubsub) InitFlags() {
	pf := ps.GetPrefix()

	flag.BoolVar(&ps.logEnabled, pf+"-log-enabled", true, "Enable logger of pubsub system")
	flag.BoolVar(&ps.gracefulStop, pf+"-graceful-stop", false, "Enable graceful shutdown")
}

func (ps *pubsub) Configure() error {
	ps.logger = logger.GetCurrent().GetLogger(ps.GetPrefix())
	return nil
}

func (ps *pubsub) Run() error {
	_ = ps.Configure()
	ps.isStopping = false

	go ps.listen()

	if ps.logEnabled {
		ps.logger.Infoln(fmt.Sprintf("started"))
	}

	return nil
}

func (ps *pubsub) Stop() <-chan bool {
	c := make(chan bool)

	go func() {
		if ps.gracefulStop {
			ps.wg.Wait()
		}

		ps.locker.Lock()

		for _, chans := range ps.mapChannel {
			for _, c := range chans {
				close(c)
			}
		}
		ps.mapChannel = make(map[pb.Channel][]chan *pb.Event)
		ps.locker.Unlock()

		if ps.logEnabled {
			ps.logger.Infoln(fmt.Sprintf("Stopped"))
		}
		c <- true
	}()

	return c
}

func (ps *pubsub) Publish(ctx context.Context, channel pb.Channel, data *pb.Event) error {
	if ps.isStopping {
		return nil
	}

	// Need to know what channel event will push to
	data.SetChannel(channel)

	go func() {
		ps.messageQueue <- data

		if ps.logEnabled {
			ps.logger.Debugln(fmt.Sprintf("new event enqueue: %s", data.String()))
		}

		ps.wg.Add(1)
	}()
	return nil
}

func (ps *pubsub) Subscribe(ctx context.Context, channel pb.Channel) (ch <-chan *pb.Event, close func()) {
	c := make(chan *pb.Event, 1)

	ps.locker.Lock()
	if m, ok := ps.mapChannel[channel]; ok {
		ps.mapChannel[channel] = []chan *pb.Event{c}
	} else {
		ps.mapChannel[channel] = append(m, c)
	}

	ps.locker.Unlock()

	if ps.logEnabled {
		ps.logger.Debugln(fmt.Sprintf("new subscriber on %s", channel))
	}

	return c, func() {
		ps.locker.Lock()
		m := ps.mapChannel[channel]

		for i := range m {
			if m[i] == c {
				ps.mapChannel[channel] = append(m[:i], m[i+1:]...)
				break
			}
		}

		ps.locker.Unlock()

		if ps.logEnabled {
			ps.logger.Debugln(fmt.Sprintf("remove a subscriber on %s", channel))
		}
	}
}

func (ps *pubsub) listen() {
	go func() {
		for {
			select {
			case <-ps.stopChan:
				ps.isStopping = true

				if ps.logEnabled {
					ps.logger.Infoln(fmt.Sprintf("stopping..."))
				}

				return
			case evt := <-ps.messageQueue:
				evt.SetAck(func() { ps.wg.Done() })

				if ps.logEnabled {
					ps.logger.Debugln(fmt.Sprintf("event did dequeue: %s", evt.String()))
				}

				ps.locker.RLock()
				chans, ok := ps.mapChannel[evt.GetChannel()]
				ps.locker.RUnlock()

				if ok {
					if len(chans) > 0 {
						ps.wg.Add(1)
					}

					for _, evtChan := range chans {
						go func(c chan *pb.Event) { c <- evt }(evtChan)
					}

				}
			}
		}
	}()
}
