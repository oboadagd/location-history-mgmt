// Package migration implements initialization of relational database
package migration

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/gommon/log"
)

// Init creates relational database if it is not existing
func Init(db *pg.DB) {
	// create a new collection with gopg_migrations table
	c := migrations.NewCollection()
	c.DisableSQLAutodiscover(true)
	err := c.DiscoverSQLMigrations(fmt.Sprintf("migration"))
	if err != nil {
		panic(err.Error())
	}

	// create gopg_migrations table if not exists
	_, _, _ = c.Run(db, "init")

	// run migration files
	oldVersion, newVersion, err := c.Run(db, "up")
	if err != nil {
		panic(err.Error())
	}
	if newVersion != oldVersion {
		log.Infof("migrated from version %d to %d", oldVersion, newVersion)
	} else {
		log.Infof("version is %d", oldVersion)
	}
}
