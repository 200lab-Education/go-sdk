/*
 * @author           Viet Tran <viettranx@gmail.com>
 * @copyright        2019 200lab <core@200lab.io>
 * @license          Apache-2.0
 */

package sdkcm

import (
	"github.com/globalsign/mgo"
	"github.com/jinzhu/gorm"
)

type GormManager interface {
	DB() *gorm.DB
}

type MgoManager interface {
	Session() *mgo.Session
}

// Remove dependent from Service Context, fix cycle import
type SC interface {
	MustGet(key string) interface{}
}

type mongo struct {
	key string
	sc  SC
}

func NewMongo(key string, sc SC) *mongo {
	return &mongo{key: key, sc: sc}
}

func (m *mongo) Session() *mgo.Session {
	return m.sc.MustGet(m.key).(*mgo.Session)
}

type sqlGorm struct {
	sc  SC
	key string
}

func NewSQLGorm(key string, sc SC) *sqlGorm {
	return &sqlGorm{key: key, sc: sc}
}

func (m *sqlGorm) DB() *gorm.DB {
	return m.sc.MustGet("mdb").(*gorm.DB)
}
