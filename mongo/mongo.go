package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Host         string
	Port         string
	UserName     string
	Password     string
	DatabaseName string

	Client   *mongo.Client
	Database *mongo.Database
}

func (m *Mongo) getConnectString() string {
	conn := "mongodb://"
	if m.UserName != "" {
		conn += fmt.Sprintf("%s:%s@", m.UserName, m.Password)
	}
	conn += fmt.Sprintf("%s:%s/%s", m.Host, m.Port, m.DatabaseName)
	return conn
}

func (m *Mongo) Connect(ctx context.Context) error {
	mgoClientOptions := options.Client().ApplyURI(m.getConnectString()).SetConnectTimeout(time.Minute).SetMaxPoolSize(5)
	// Connect to MongoDB
	var err error
	if m.Client, err = mongo.Connect(ctx, mgoClientOptions); err != nil {
		return err
	}

	m.Database = m.Client.Database(m.DatabaseName)
	return nil
}

func (m *Mongo) Disconnect(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
