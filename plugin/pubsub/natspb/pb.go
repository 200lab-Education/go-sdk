package natspb

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/200Lab-Education/go-sdk/logger"
	pb "github.com/200Lab-Education/go-sdk/plugin/pubsub"
	"github.com/nats-io/nats.go"
)

type NatsOpt struct {
	prefix   string
	server   string
	username string
	password string
	token    string
}

type natspb struct {
	name      string
	logger    logger.Logger
	nc        *nats.Conn
	isRunning bool
	*NatsOpt
}

func NewNatsPubSub(name string, prefix string) *natspb {
	return &natspb{
		name: name,
		NatsOpt: &NatsOpt{
			prefix: prefix,
		},
		isRunning: false,
	}
}

func (n *natspb) GetPrefix() string {
	if n.prefix == "" {
		return n.name
	}
	return n.prefix
}

func (n *natspb) Get() interface{} {
	return n
}

func (n *natspb) Name() string {
	return n.name
}

func (n *natspb) InitFlags() {
	prefix := n.prefix
	if n.prefix != "" {
		prefix += "-"
	}

	flag.StringVar(&n.server, prefix+"nats-server", "", "Nats connect server. Ex: \"nats://..., nats://\"")
	flag.StringVar(&n.username, prefix+"nats-username", "", "Nats username")
	flag.StringVar(&n.password, prefix+"nats-password", "", "Nats password")
	flag.StringVar(&n.token, prefix+"nats-token", "", "Nats token")
}

func (n *natspb) Configure() error {
	if n.isRunning {
		return nil
	}
	n.logger = logger.GetCurrent().GetLogger(n.name)
	n.logger.Info("Connecting to Nats at ", n.server, " ...")

	var options []nats.Option

	options = append(options,
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			n.logger.Errorf("Got disconnected! Reason: %q\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			n.logger.Errorf("Got reconnected to %v!\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			n.logger.Errorf("Connection closed. Reason: %q\n", nc.LastError())
		}))

	if n.username != "" {
		options = append(options, nats.UserInfo(n.username, n.password))
	}
	if n.token != "" {
		options = append(options, nats.Token(n.token))
	}

	nc, err := nats.Connect(n.server,
		options...)

	if err != nil {
		return err
	}

	n.nc = nc
	n.isRunning = true
	return nil
}

func (n *natspb) Run() error {
	return n.Configure()
}

func (n *natspb) Stop() <-chan bool {
	if n.nc != nil {

		err := n.nc.Drain()
		if err != nil {
			n.logger.Errorf("Error when drain nats connection: %q\n", err)
		}
	}
	n.isRunning = false

	c := make(chan bool)
	go func() { c <- true }()
	return c
}

func (n *natspb) Publish(ctx context.Context, channel pb.Channel, data *pb.Event) error {
	dataByte, err := json.Marshal(data.Data)

	if err != nil {
		n.logger.Errorln(err)
		return err
	}

	if err := n.nc.Publish(string(channel), dataByte); err != nil {
		n.logger.Errorln(err)
		return err
	}

	//if err := n.nc.Flush(); err != nil {
	//	n.logger.Errorln(err)
	//	return err
	//}

	return nil
}

func (n *natspb) Subscribe(ctx context.Context, channel pb.Channel) (c <-chan *pb.Event, cl func()) {
	ch := make(chan *pb.Event)

	sub, err := n.nc.Subscribe(string(channel), func(msg *nats.Msg) {
		evt := &pb.Event{
			Channel:    channel,
			RemoteData: msg.Data,
		}

		ch <- evt
	})

	if err != nil {
		n.logger.Errorln(err)
	}

	return ch, func() {
		_ = sub.Unsubscribe()
		close(ch)
	}
}
