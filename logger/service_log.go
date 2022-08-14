package logger

import (
	"flag"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var (
	DefaultStdLogger = NewAppLogService(&Config{
		BasePrefix:   "core",
		DefaultLevel: "trace",
	})
)

type Config struct {
	DefaultLevel string
	BasePrefix   string
}

type ServiceLogger interface {
	GetLogger(prefix string) Logger
}

// A default app logger
// Just write everything to console
type stdLogger struct {
	logger   *logrus.Logger
	cfg      Config
	logLevel string
}

func NewAppLogService(config *Config) *stdLogger {
	if config == nil {
		config = &Config{}
	}

	//flag.StringVar(&config.DefaultLevel, "log-level", "info", "Log level: panic | fatal | error | warn | info | debug | trace")

	if config.DefaultLevel == "" {
		config.DefaultLevel = "info"
	}

	logger := logrus.New()
	logger.Formatter = logrus.Formatter(&prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	return &stdLogger{
		logger:   logger,
		cfg:      *config,
		logLevel: config.DefaultLevel,
	}
}

func (s *stdLogger) GetLogger(prefix string) Logger {
	var entry *logrus.Entry

	prefix = s.cfg.BasePrefix + "." + prefix
	prefix = strings.Trim(prefix, ".")
	if prefix == "" {
		entry = logrus.NewEntry(s.logger)
	} else {
		entry = s.logger.WithField("prefix", prefix)
	}

	l := &logger{entry}
	var log Logger = l

	return log
}

// Implement Runnable interface
func (s *stdLogger) Name() string { return "file-logger" }
func (s *stdLogger) InitFlags() {
	flag.StringVar(&s.logLevel, "log-level", s.cfg.DefaultLevel, "Log level: panic | fatal | error | warn | info | debug | trace")
}
func (s *stdLogger) Configure() error {
	lv := mustParseLevel(s.logLevel)
	s.logger.SetLevel(lv)
	return nil
}

func (s *stdLogger) Run() error { return s.Configure() }
func (s *stdLogger) Stop() <-chan bool {
	c := make(chan bool)
	go func() { c <- true }()

	return c
}
