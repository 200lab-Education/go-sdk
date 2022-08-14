// Copyright (c) 2019, Viet Tran, 200Lab Team.

package goservice

import (
	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/gin-gonic/gin"
)

// Convenience option method for creating/initializing a service
type Option func(*service)

// HTTP Server Handler for register some routes and gin handlers
type HttpServerHandler = func(*gin.Engine)

// A kind of server job
type Function func(ServiceContext) error

// The storage store all db connection in service
type Storage interface {
	Get(prefix string) (interface{}, bool)
	MustGet(prefix string) interface{}
}

type PrefixRunnable interface {
	HasPrefix
	Runnable
}

type HasPrefix interface {
	GetPrefix() string
	Get() interface{}
}

// The heart of SDK, Service represents for a real micro service
// with its all components
type Service interface {
	// A part of Service, it's passed to all handlers/functions
	ServiceContext
	// Name of the service
	Name() string
	// Version of the service
	Version() string
	// Gin HTTP Server wrapper
	HTTPServer() HttpServer
	// Init with options, they can be db connections or
	// anything the service need handle before starting
	Init() error
	// This method returns service if it is registered on discovery
	IsRegistered() bool
	// Start service and its all component.
	// It will be stopped if any service return error
	Start() error
	// Stop service and its all component.
	Stop()
	// Method export all flags to std/terminal
	// We might use: "> .env" to move its content .env file
	OutEnv()
}

// Service Context: A wrapper for all things needed for developing a service
type ServiceContext interface {
	// Logger for a specific service, usually it has a prefix to distinguish
	// with each others
	Logger(prefix string) logger.Logger
	// Get component with prefix
	Get(prefix string) (interface{}, bool)
	MustGet(prefix string) interface{}
	Env() string
}

// Runnable is an abstract object in SDK
// Almost components are Runnable. SDK will manage their lifecycle
// InitFlags -> Configure -> Run -> Stop
type Runnable interface {
	Name() string
	InitFlags()
	Configure() error
	Run() error
	Stop() <-chan bool
}

// GIN HTTP server for REST API
type HttpServer interface {
	Runnable
	// Add handlers to GIN
	AddHandler(HttpServerHandler)
	// Return server config
	//GetConfig() http_server.Config
	// URI that the server is listening
	URI() string
}
