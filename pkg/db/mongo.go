package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Intializing MongoDb and connecting to its instance

func InitMongo(uri string) (*mongo.Client,error){
	ctx,cancel :=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	// create a new client and connect to mongoDB
	client,err :=mongo.Connect(ctx,options.Client().ApplyURI(uri))
	if err != nil{
		log.Fatal("Failed to connect to mongoDB",err)
	}

	// check the connection

	err =client.Ping(ctx,nil)
	if err != nil{
		log.Fatal("Failed to ping to mongoDB",err)
	}

	fmt.Println("Connected to mongoDB")
	return client,nil
} 

func GetCollection(client *mongo.Client,dbName,collectionName string)*mongo.Collection{
	return client.Database(dbName).Collection(collectionName)
}

