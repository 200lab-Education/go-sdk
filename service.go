// Copyright (c) 2019, Viet Tran, 200Lab Team.

package goservice

import (
	"flag"
	"fmt"
	"github.com/200Lab-Education/go-sdk/httpserver"
	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	DevEnv     = "dev"
	StgEnv     = "stg"
	PrdEnv     = "prd"
	DefaultEnv = DevEnv
)

type service struct {
	name         string
	version      string
	env          string
	opts         []Option
	subServices  []Runnable
	initServices map[string]PrefixRunnable
	isRegister   bool
	logger       logger.Logger
	httpServer   HttpServer
	signalChan   chan os.Signal
	cmdLine      *AppFlagSet
	stopFunc     func()
}

func New(opts ...Option) Service {
	sv := &service{
		opts:         opts,
		signalChan:   make(chan os.Signal, 1),
		subServices:  []Runnable{},
		initServices: map[string]PrefixRunnable{},
	}

	// init default logger
	logger.InitServLogger(false)
	sv.logger = logger.GetCurrent().GetLogger("service")

	for _, opt := range opts {
		opt(sv)
	}

	//// Http server
	httpServer := httpserver.New(sv.name)
	sv.httpServer = httpServer

	sv.subServices = append(sv.subServices, httpServer)

	sv.initFlags()

	if sv.name == "" {
		if len(os.Args) >= 2 {
			sv.name = strings.Join(os.Args[:2], " ")
		}
	}

	loggerRunnable := logger.GetCurrent().(Runnable)
	loggerRunnable.InitFlags()

	sv.cmdLine = newFlagSet(sv.name, flag.CommandLine)
	sv.parseFlags()

	_ = loggerRunnable.Configure()

	return sv
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Version() string {
	return s.version
}

func (s *service) Init() error {
	for _, dbSv := range s.initServices {
		if err := dbSv.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) IsRegistered() bool {
	return s.isRegister
}

func (s *service) Start() error {
	signal.Notify(s.signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	c := s.run()
	//s.stopFunc = s.activeRegistry()

	for {
		select {
		case err := <-c:
			if err != nil {
				s.logger.Error(err.Error())
				s.Stop()
				return err
			}

		case sig := <-s.signalChan:
			s.logger.Infoln(sig)
			switch sig {
			case syscall.SIGHUP:
				return nil
			default:
				s.Stop()
				return nil
			}
		}
	}
}

func (s *service) initFlags() {
	flag.StringVar(&s.env, "app-env", DevEnv, "Env for service. Ex: dev | stg | prd")

	for _, subService := range s.subServices {
		subService.InitFlags()
	}

	for _, dbService := range s.initServices {
		dbService.InitFlags()
	}
}

// Run service and its components at the same time
func (s *service) run() <-chan error {
	c := make(chan error, 1)

	// Start all services
	for _, subService := range s.subServices {
		go func(subSv Runnable) { c <- subSv.Run() }(subService)
	}

	return c
}

// Stop service and stop its components at the same time
func (s *service) Stop() {
	s.logger.Infoln("Stopping service...")
	stopChan := make(chan bool)
	for _, subService := range s.subServices {
		go func(subSv Runnable) { stopChan <- <-subSv.Stop() }(subService)
	}

	for _, dbSv := range s.initServices {
		go func(subSv Runnable) { stopChan <- <-subSv.Stop() }(dbSv)
	}

	for i := 0; i < len(s.subServices)+len(s.initServices); i++ {
		<-stopChan
	}

	//s.stopFunc()
	s.logger.Infoln("service stopped")
}

func (s *service) RunFunction(fn Function) error {
	return fn(s)
}

func (s *service) HTTPServer() HttpServer {
	return s.httpServer
}

func (s *service) Logger(prefix string) logger.Logger {
	return logger.GetCurrent().GetLogger(prefix)
}

func (s *service) OutEnv() {
	s.cmdLine.GetSampleEnvs()
}

func (s *service) parseFlags() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_, err := os.Stat(envFile)
	if err == nil {
		err := godotenv.Load(envFile)
		if err != nil {
			s.logger.Fatalf("Loading env(%s): %s", envFile, err.Error())
		}
	} else if envFile != ".env" {
		s.logger.Fatalf("Loading env(%s): %s", envFile, err.Error())
	}

	s.cmdLine.Parse([]string{})
}

// Service must have a name for service discovery and logging/monitoring
func WithName(name string) Option {
	return func(s *service) { s.name = name }
}

// Every deployment needs a specific version
func WithVersion(version string) Option {
	return func(s *service) { s.version = version }
}

// Service will write log data to file with this option
func WithFileLogger() Option {
	return func(s *service) {
		logger.InitServLogger(true)
	}
}

// Add Runnable component to SDK
// These components will run parallel in when service run
func WithRunnable(r Runnable) Option {
	return func(s *service) { s.subServices = append(s.subServices, r) }
}

// Add init component to SDK
// These components will run sequentially before service run
func WithInitRunnable(r PrefixRunnable) Option {
	return func(s *service) {
		if _, ok := s.initServices[r.GetPrefix()]; ok {
			log.Fatal(fmt.Sprintf("prefix %s is duplicated", r.GetPrefix()))
		}

		s.initServices[r.GetPrefix()] = r
	}
}

func (s *service) Get(prefix string) (interface{}, bool) {
	is, ok := s.initServices[prefix]

	if !ok {
		return nil, ok
	}

	return is.Get(), true
}

func (s *service) MustGet(prefix string) interface{} {
	db, ok := s.Get(prefix)

	if !ok {
		panic(fmt.Sprintf("can not get %s\n", prefix))
	}

	return db
}

func (s *service) Env() string { return s.env }
