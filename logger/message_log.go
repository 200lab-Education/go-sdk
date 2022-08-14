package logger

import (
	"flag"
	"github.com/sirupsen/logrus"
)

var (
	DefaultMessageLogger = NewMessageLogService(&Config{
		BasePrefix:   "core",
		DefaultLevel: "info",
	})
)

// Used to log message that need send to remote such as Elastic Search
type messageLogger struct {
	*stdLogger

	// log error message
	log Logger
	// flags
	logPath string
}

func NewMessageLogService(config *Config) *messageLogger {
	if config == nil {
		config = &Config{}
	}

	appLog := NewAppLogService(config)
	appLog.logger.Formatter = &logrus.JSONFormatter{}

	newLog := logrus.New()
	newLog.Formatter = logrus.Formatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	log := &logger{logrus.NewEntry(newLog)}

	return &messageLogger{
		stdLogger: appLog,
		log:       log,
	}
}

// Implement Runnable interface
func (m *messageLogger) Name() string { return "file-logger" }

func (m *messageLogger) InitFlags() {
	flag.StringVar(&m.logPath, "logfile", "", "file to write log to. Default write to console")
	m.stdLogger.InitFlags()
}

func (m *messageLogger) Configure() error {

	if m.logPath == "" {
		lv := mustParseLevel(m.stdLogger.cfg.DefaultLevel)
		m.stdLogger.logger.SetLevel(lv)
		return nil
	}

	out, err := newReloadFile(m.logPath)
	if err != nil {
		m.log.Fatal("Fail to open log file: ", err.Error())
	}
	m.logger.Out = out

	return nil
}

func (m *messageLogger) Run() error {
	return m.Configure()
}

func (m *messageLogger) Stop() <-chan bool {
	c := make(chan bool)

	go func() {
		if file, ok := m.logger.Out.(*reloadFile); ok {
			file.Close()
		}
		c <- true
	}()

	return c
}
