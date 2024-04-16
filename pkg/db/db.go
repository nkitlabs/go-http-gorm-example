package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nkitlabs/go-http-gorm-example/pkg/config"
)

// Init initializes the database connection.
func Init(conf config.DBConn) (*gorm.DB, error) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName,
	)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
