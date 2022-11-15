package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AnimeTrackerr/backend/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


func ConnectDB(URI string) utils.ClientInfo {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	var err1 error = client.Connect(ctx)
	var err2 error = client.Ping(ctx, readpref.Primary())

	if err1 != nil{
		log.Fatal(err1)
	} else if err2 != nil{
		log.Fatal(err2)
	} else {
		fmt.Print("Successfully connected to DB\n")
	}
	
	return utils.ClientInfo {
		Client: client,
		Ctx: ctx,
	}
}

func GetCollection(client *mongo.Client, dbName string, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}
