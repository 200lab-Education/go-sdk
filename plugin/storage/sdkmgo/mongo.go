/*
 * @author          Viet Tran <viettranx@gmail.com>
 * @copyright       2019 Viet Tran <viettranx@gmail.com>
 * @license         Apache-2.0
 */

package sdkmgo

import (
	"flag"
	"github.com/200Lab-Education/go-sdk/logger"
	"math"
	"sync"
	"time"

	"github.com/globalsign/mgo"
)

var (
	defaultDBName  = "defaultMongoDB"
	DefaultMongoDB = getDefaultMongoDB()
)

const retryCount = 10

type MongoDBOpt struct {
	MgoUri       string
	Prefix       string
	PingInterval int // in seconds
}

type mongoDB struct {
	name      string
	logger    logger.Logger
	session   *mgo.Session
	isRunning bool
	once      *sync.Once
	*MongoDBOpt
}

func getDefaultMongoDB() *mongoDB {
	return NewMongoDB(defaultDBName, "")
}

func NewMongoDB(name, prefix string) *mongoDB {
	return &mongoDB{
		MongoDBOpt: &MongoDBOpt{
			Prefix: prefix,
		},
		name:      name,
		isRunning: false,
		once:      new(sync.Once),
	}
}

func (mgDB *mongoDB) GetPrefix() string {
	return mgDB.Prefix
}

func (mgDB *mongoDB) Name() string {
	return mgDB.name
}

func (mgDB *mongoDB) InitFlags() {
	prefix := mgDB.Prefix
	if mgDB.Prefix != "" {
		prefix += "-"
	}

	flag.StringVar(&mgDB.MgoUri, prefix+"mgo-uri", "", "MongoDB connection-string. Ex: mongodb://...")
	flag.IntVar(&mgDB.PingInterval, prefix+"mgo-ping-interval", 5, "MongoDB ping check interval")
}

func (mgDB *mongoDB) isDisabled() bool {
	return mgDB.MgoUri == ""
}

func (mgDB *mongoDB) Configure() error {
	if mgDB.isDisabled() || mgDB.isRunning {
		return nil
	}

	mgDB.logger = logger.GetCurrent().GetLogger(mgDB.name)
	mgDB.logger.Info("Connect to Mongodb at ", mgDB.MgoUri, " ...")

	var err error
	mgDB.session, err = mgDB.getConnWithRetry(retryCount)
	if err != nil {
		mgDB.logger.Error("Error connect to mongodb at ", mgDB.MgoUri, ". ", err.Error())
		return err
	}
	mgDB.isRunning = true
	return nil
}

func (mgDB *mongoDB) Cleanup() {
	if mgDB.isDisabled() {
		return
	}

	if mgDB.session != nil {
		mgDB.session.Close()
	}
}

func (mgDB *mongoDB) Run() error {
	return mgDB.Configure()
}

func (mgDB *mongoDB) Stop() <-chan bool {
	if mgDB.session != nil {
		mgDB.session.Close()
	}
	mgDB.isRunning = false

	c := make(chan bool)
	go func() { c <- true }()
	return c
}

func (mgDB *mongoDB) Get() interface{} {
	mgDB.once.Do(func() {
		if !mgDB.isRunning && !mgDB.isDisabled() {
			if db, err := mgDB.getConnWithRetry(math.MaxInt32); err == nil {
				mgDB.session = db
				mgDB.isRunning = true
			} else {
				mgDB.logger.Fatalf("%s connection cannot reconnect\n", mgDB.name)
			}
		}
	})

	if mgDB.session == nil {
		return nil
	}
	return mgDB.session.New()
}

func (mgDB *mongoDB) getConnWithRetry(retryCount int) (*mgo.Session, error) {
	db, err := mgo.Dial(mgDB.MgoUri)

	if err != nil {
		for {
			time.Sleep(time.Second * 1)
			mgDB.logger.Errorf("Retry to connect %s.\n", mgDB.name)
			db, err = mgo.Dial(mgDB.MgoUri)

			if err == nil {
				go mgDB.reconnectIfNeeded()
				break
			}
		}
	} else {
		go mgDB.reconnectIfNeeded()
	}

	return db, err
}

func (mgDB *mongoDB) reconnectIfNeeded() {
	conn := mgDB.session
	for {
		if err := conn.Ping(); err != nil {
			conn.Close()
			mgDB.logger.Errorf("%s connection is gone, try to reconnect\n", mgDB.name)
			mgDB.isRunning = false
			mgDB.once = new(sync.Once)

			mgDB.Get().(*mgo.Session).Close()
			return
		}
		time.Sleep(time.Second * time.Duration(mgDB.PingInterval))
	}
}
