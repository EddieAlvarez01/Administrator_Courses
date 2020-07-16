package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/EddieAlvarez01/administrator_courses/utilities"
)

//GetConnection get connection of mongodb
func GetConnection() *mongo.Client {
	config, err := utilities.GetDatabaseConfiguration()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	client, err2 := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", config.Server, config.Port)))
	if err2 != nil {
		log.Fatal(err2)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	//PING TO DB
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return client
}
