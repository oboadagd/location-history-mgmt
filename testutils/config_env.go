// Package testutils implements utils for testing of microservice location-history-mgmt.
package testutils

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/kelseyhightower/envconfig"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-common/recordtype"
	"strings"
)

// pgOptions returns database configuration options
func pgOptions() *pg.Options {

	envconfig.Process("LIST", &recordtype.Cfg)

	return &pg.Options{
		Addr:     fmt.Sprintf("%s:%d", recordtype.Cfg.DBHost, recordtype.Cfg.DBPort),
		User:     recordtype.Cfg.DBUser,
		Password: recordtype.Cfg.DBPass,
	}
}

// GetTestDB returns a database instance connection
func GetTestDB() *pg.DB {
	return pg.Connect(pgOptions())
}

// CreateSchema Schema in the mock DB still needs to be created
func CreateSchema(db *pg.DB) error {

	err := db.Model((*dto.Location)(nil)).CreateTable(&orm.CreateTableOptions{
		// set Temp=True so no tables/data are actually created
		Temp: true,
	})
	if err != nil {
		return err
	}

	err = db.Model((*dto.LocationHistory)(nil)).CreateTable(&orm.CreateTableOptions{
		// set Temp=True so no tables/data are actually created
		Temp: true,
	})
	if err != nil {
		return err
	}

	return nil
}

// DropSchema Schema in the mock DB still needs to be dropped
func DropSchema(db *pg.DB) error {

	err := db.Model((*dto.Location)(nil)).DropTable(nil)
	if err != nil {
		return err
	}

	return nil
}

// GetLocation returns an instanced *dto.SaveLocationRequest
func GetLocation() *dto.SaveLocationRequest {
	return &dto.SaveLocationRequest{
		UserName:  "usernamesample",
		Latitude:  10,
		Longitude: 10,
	}
}

// GetLocation returns an instanced *dto.CreateLocationHistoryRequest
func GetLocationHistory() *dto.CreateLocationHistoryRequest {
	return &dto.CreateLocationHistoryRequest{
		UserName:  "usernamesample",
		Latitude:  10,
		Longitude: 10,
		Distance:  1,
	}
}

// EvaluateErrConditions returns evaluation of values. Returns false if at least one evaluation is
// not met, returns true otherwise
func EvaluateErrConditions(errMsg string, values []string) bool {

	for _, value := range values {
		if !strings.Contains(strings.ToLower(errMsg), value) {
			return true
		}
	}

	return false
}
