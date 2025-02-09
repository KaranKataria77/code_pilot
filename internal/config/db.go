package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Client

func ConnectDB() {
	uri := GetEnv("MONGO_URL", "mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error connecting Database ", err)
	}
	MongoDB = client
	log.Println("MongoDB connected")
}

func GetMongoCollection(name string) *mongo.Collection {
	return MongoDB.Database("code_pilot").Collection(name)
}
