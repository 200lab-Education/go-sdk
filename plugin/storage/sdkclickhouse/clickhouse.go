/*
 * @author           Viet Tran <viettranx@gmail.com>
 * @copyright        2020 200lab <core@200lab.io>
 * @license          Apache-2.0
 */

package sdkclickhouse

import (
	"flag"
	"github.com/200Lab-Education/go-sdk/logger"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

const retryCount = 10

type CHDBOpt struct {
	ChUri        string
	Prefix       string
	PingInterval int // in seconds
}

type clickhouseDB struct {
	name      string
	logger    logger.Logger
	session   *sqlx.DB
	isRunning bool
	once      *sync.Once
	*CHDBOpt
}

func NewClickHouseDB(name, prefix string) *clickhouseDB {
	return &clickhouseDB{
		CHDBOpt: &CHDBOpt{
			Prefix: prefix,
		},
		name:      name,
		isRunning: false,
		once:      new(sync.Once),
	}
}

func (chDB *clickhouseDB) GetPrefix() string {
	return chDB.Prefix
}

func (chDB *clickhouseDB) Name() string {
	return chDB.name
}

func (chDB *clickhouseDB) InitFlags() {
	prefix := chDB.Prefix
	if chDB.Prefix != "" {
		prefix += "-"
	}

	flag.StringVar(&chDB.ChUri, prefix+"clickhouse-uri", "", "ClickHouse connection-string. Ex: tcp://host1:9000?username=user&password=qwerty&database=clicks")
	flag.IntVar(&chDB.PingInterval, prefix+"clickhouse-ping-interval", 5, "ClickHouse ping check interval")
}

func (chDB *clickhouseDB) isDisabled() bool {
	return chDB.ChUri == ""
}

func (chDB *clickhouseDB) Configure() error {
	if chDB.isDisabled() || chDB.isRunning {
		return nil
	}

	chDB.logger = logger.GetCurrent().GetLogger(chDB.name)
	chDB.logger.Info("Connect to ClickHouse at ", chDB.ChUri, " ...")

	var err error
	chDB.session, err = chDB.getConnWithRetry()
	if err != nil {
		chDB.logger.Error("Error connect to ClickHouse at ", chDB.ChUri, ". ", err.Error())
		return err
	}
	chDB.isRunning = true
	return nil
}

func (chDB *clickhouseDB) Cleanup() {
	if chDB.isDisabled() {
		return
	}

	if chDB.session != nil {
		_ = chDB.session.Close()
	}
}

func (chDB *clickhouseDB) Run() error {
	return chDB.Configure()
}

func (chDB *clickhouseDB) Stop() <-chan bool {
	if chDB.session != nil {
		_ = chDB.session.Close()
	}
	chDB.isRunning = false

	c := make(chan bool)
	go func() { c <- true }()
	return c
}

func (chDB *clickhouseDB) Get() interface{} {
	chDB.once.Do(func() {
		if !chDB.isRunning && !chDB.isDisabled() {
			if db, err := chDB.getConnWithRetry(); err == nil {
				chDB.session = db
				chDB.isRunning = true
			} else {
				chDB.logger.Fatalf("%s connection cannot reconnect\n", chDB.name)
			}
		}
	})

	if chDB.session == nil {
		return nil
	}
	return chDB.session
}

func (chDB *clickhouseDB) getConnWithRetry() (*sqlx.DB, error) {
	db, err := sqlx.Connect("clickhouse", chDB.ChUri)

	if err != nil {
		for {
			time.Sleep(time.Second * 1)
			chDB.logger.Errorf("Retry to connect %s.\n", chDB.name)
			db, err = sqlx.Connect("clickhouse", chDB.ChUri)

			if err == nil {
				go chDB.reconnectIfNeeded()
				break
			}
		}
	} else {
		go chDB.reconnectIfNeeded()
	}

	return db, err
}

func (chDB *clickhouseDB) reconnectIfNeeded() {
	conn := chDB.session
	for {
		if err := conn.Ping(); err != nil {
			_ = conn.Close()
			chDB.logger.Errorf("%s connection is gone, try to reconnect\n", chDB.name)
			chDB.isRunning = false
			chDB.once = new(sync.Once)

			_ = chDB.Get().(*sqlx.DB).Close()
			return
		}
		time.Sleep(time.Second * time.Duration(chDB.PingInterval))
	}
}
