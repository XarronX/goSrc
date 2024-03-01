package db

import (
	"fmt"
	"sync"

	// _ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "eqconsultancy"
)

var connection *gorm.DB
var connIssue error
var once sync.Once

func init() {
	once.Do(func() {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		connection, connIssue = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	})
}

func Connect() (*gorm.DB, error) {
	if connIssue != nil {
		return nil, connIssue
	}
	return connection, nil
}
