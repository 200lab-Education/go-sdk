package sckio

import (
	"flag"
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/200Lab-Education/go-sdk/sdkcm"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
)

type Socket interface {
	Id() string
	Rooms() []string
	Request() *http.Request
	On(event string, f interface{}) error
	Emit(event string, args ...interface{}) error
	Join(room string) error
	Leave(room string) error
	Disconnect()
	BroadcastTo(room, event string, args ...interface{}) error
}

type AppSocket interface {
	ServiceContext() goservice.ServiceContext
	Logger() logger.Logger
	CurrentUser() sdkcm.Requester
	SetCurrentUser(sdkcm.Requester)
	BroadcastToRoom(room, event string, args ...interface{})
	String() string
	Socket
}

type Config struct {
	Name          string
	MaxConnection int
}

type sckServer struct {
	Config
	io     *socketio.Server
	logger logger.Logger
}

func New(name string) *sckServer {
	return &sckServer{
		Config: Config{Name: name},
	}
}

type ObserverProvider interface {
	AddObservers(server *socketio.Server, sc goservice.ServiceContext, l logger.Logger) func(socketio.Socket)
}

func (s *sckServer) StartRealtimeServer(engine *gin.Engine, sc goservice.ServiceContext, op ObserverProvider) {
	server, err := socketio.NewServer([]string{"websocket"})
	if err != nil {
		log.Fatal(err)
	}

	server.SetMaxConnection(s.MaxConnection)
	s.io = server

	_ = s.io.On("connection", op.AddObservers(server, sc, s.logger))

	engine.GET("/socket.io/", gin.WrapH(server))
	engine.POST("/socket.io/", gin.WrapH(server))
}

func (s *sckServer) GetPrefix() string {
	return s.Config.Name
}

func (s *sckServer) Get() interface{} {
	return s
}

func (s *sckServer) Name() string {
	return s.Config.Name
}

func (s *sckServer) InitFlags() {
	pre := s.GetPrefix()
	flag.IntVar(&s.MaxConnection, fmt.Sprintf("%s-max-connection", pre), 2000, "socket max connection")
}

func (s *sckServer) Configure() error {
	s.logger = logger.GetCurrent().GetLogger("io.socket")
	return nil
}

func (s *sckServer) Run() error {
	return s.Configure()
}

func (s *sckServer) Stop() <-chan bool {
	c := make(chan bool)
	go func() { c <- true }()
	return c
}
