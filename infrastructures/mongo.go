package infrastructures

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoClient(uri, username, password string) (*mongo.Client, error) {

	credential := options.Credential{
		Username: username,
		Password: password,
	}
	clientOpts := options.Client().ApplyURI(uri).SetAuth(credential)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOpts)

	err = client.Ping(ctx, nil)

	return client, err
}
