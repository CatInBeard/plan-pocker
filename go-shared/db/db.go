package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type DbClient struct {
	primaryDB *sql.DB
	slaveDB   *sql.DB
}

var (
	instance   *DbClient
	once       sync.Once
	errReplica error
)

func GetDbClient() (*DbClient, error) {
	once.Do(func() {
		instance = &DbClient{}
		var errPrimary error
		errPrimary, errReplica = instance.init()
		if errPrimary != nil {
			panic(errPrimary)
		}

	})
	return instance, nil
}

func (d *DbClient) init() (error, error) {
	primaryDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_PRIMARY_USER"),
		os.Getenv("DB_PRIMARY_PASSWORD"),
		os.Getenv("DB_PRIMARY_HOST"),
		os.Getenv("DB_PRIMARY_PORT"),
		os.Getenv("DB_PRIMARY_NAME"),
	)

	replicaDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_REPLICA_USER"),
		os.Getenv("DB_REPLICA_PASSWORD"),
		os.Getenv("DB_REPLICA_HOST"),
		os.Getenv("DB_REPLICA_PORT"),
		os.Getenv("DB_REPLICA_NAME"),
	)

	var err error
	d.primaryDB, err = sql.Open("mysql", primaryDSN)
	if err != nil {
		return err, nil
	}

	d.slaveDB, err = sql.Open("mysql", replicaDSN)

	if err != nil {
		return nil, err
	}

	d.primaryDB.SetMaxOpenConns(500)
	d.slaveDB.SetMaxOpenConns(1000)

	return nil, nil
}

func (d *DbClient) ExecuteUpdate(query string, args ...interface{}) (sql.Result, error) {
	return d.primaryDB.Exec(query, args...)
}

func (d *DbClient) ExecuteReadPrimary(query string, args ...interface{}) (*sql.Rows, error) {
	return d.primaryDB.Query(query, args...)
}

func (d *DbClient) ExecuteRead(query string, args ...interface{}) (*sql.Rows, error) {
	return d.slaveDB.Query(query, args...)
}
