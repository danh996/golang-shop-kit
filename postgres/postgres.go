package postgres

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresClient struct {
	DataSource string
	Database   *gorm.DB
}

func NewPostgresClient(dataSource string) (*PostgresClient, error) {

	pConfig := postgres.Config{
		DSN:                  dataSource,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}
	gConfig := &gorm.Config{
		SkipDefaultTransaction: true,
	}
	db, err := gorm.Open(postgres.New(pConfig), gConfig)

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &PostgresClient{
		DataSource: dataSource,
		Database:   db,
	}, nil
}
