/*
 * @author           Viet Tran <viettranx@gmail.com>
 * @copyright       2019 Viet Tran <viettranx@gmail.com>
 * @license           Apache-2.0
 */

package dbmigration

import (
	"fmt"
	"github.com/200Lab-Education/go-sdk/logger"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	defaultTbName = "configs"
)

type sqlMigration struct {
	db         *gorm.DB
	env        string
	tbName     string
	actVersion int
	sqlFolder  string
	logger     logger.Logger
}

func NewSQLMigration(db *gorm.DB, folder string, logger logger.Logger, opts ...Opt) *sqlMigration {
	if strings.TrimSpace(folder) == "" {
		log.Fatalln("folder cannot be empty")
	}

	if db == nil {
		log.Fatalln("db connection must not be nil")
	}

	m := &sqlMigration{db: db, sqlFolder: folder, logger: logger}
	for _, o := range opts {
		o(m)
	}

	return m
}

func (sql *sqlMigration) Name() string {
	return "SQL Migration"
}

func (sql *sqlMigration) actualVersion() int {
	var data struct {
		Value string `gorm:"value"`
	}

	if err := sql.db.Table(sql.getTableName()).
		Where("name = ?", "DB_VERSION").First(&data).Error; err != nil {
		return 0
	}

	sql.actVersion, _ = strconv.Atoi(data.Value)
	return sql.actVersion
}

func (sql *sqlMigration) Migrate() error {
	files, err := ioutil.ReadDir(sql.sqlFolder)
	if err != nil {
		return err
	}

	var versions []int
	maxVersion := 0

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".sql" {
			if n, _ := strconv.Atoi(strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))); n > 0 {
				versions = append(versions, n)

				if maxVersion < n {
					maxVersion = n
				}
			}
		}
	}

	sort.Ints(versions)

	actualVer := sql.actualVersion()
	if actualVer < maxVersion && len(versions) > 0 {
		sql.logger.Infoln("migrating sql db...")
		for _, v := range versions {
			if v <= actualVer {
				continue
			}

			mainPath, err := filepath.Abs(fmt.Sprintf("%s/%d.sql", sql.sqlFolder, v))
			subPath, err := filepath.Abs(fmt.Sprintf("%s/%d.%s.sql", sql.sqlFolder, v, sql.env))

			if err != nil {
				sql.logger.Fatalln(err)
			}

			sqlData, err := ioutil.ReadFile(mainPath)

			if err != nil {
				sql.logger.Fatalln(err)
			}

			sqlSubData, _ := ioutil.ReadFile(subPath)

			if len(sqlData) == 0 && len(sqlSubData) == 0 {
				continue
			}

			sql.logger.Infoln("migrating sql db... version:", v)

			// Because GORM can't exec multiple commands,
			// so we split it to many commands and run each
			sqlCmds := strings.Split(string(sqlData), ";")
			sqlSubCmds := strings.Split(string(sqlSubData), ";")

			if len(sqlSubCmds) > 1 {
				sqlCmds = append(sqlCmds, sqlSubCmds...)
			}

			for _, command := range sqlCmds {
				if strings.TrimSpace(command) == "" {
					continue
				}

				db := sql.db
				if err := db.Exec(command).Error; err != nil {
					sql.logger.Fatalln(err)
				}
			}

			if err := sql.db.Table(sql.getTableName()).
				Where("name = ?", "DB_VERSION").
				Updates(map[string]interface{}{"value": v}).Error; err != nil {
				sql.logger.Fatalln(err)
			}
		}

		sql.logger.Infoln("migrating sql db... done.")
	}

	return nil
}

func (sql *sqlMigration) getTableName() string {
	if sql.tbName == "" {
		return defaultTbName
	}
	return sql.tbName
}

type Opt func(*sqlMigration)

func WithTableName(s string) Opt {
	return func(sql *sqlMigration) {
		sql.tbName = s
	}
}

func WithEnv(env string) Opt {
	return func(sql *sqlMigration) {
		sql.env = env
	}
}
