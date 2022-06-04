package database

import (
	"context"
	//"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var MongoDB *mongo.Client

func MongoInit() {
	port := viper.GetString("MONGODB_PORT")
	user := viper.GetString("MONGO_INITDB_ROOT_USERNAME")
	pwd := viper.GetString("MONGO_INITDB_ROOT_PASSWORD")
	//db := viper.GetString("MONGO_INITDB_DATABASE")

	uri := "mongodb://" + user + ":" + pwd + "@localhost:" + port
	ctx := context.Background()
	MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Mongo Connect err : " + err.Error())
	}

	err = MongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Mongo Ping err : " + err.Error())
	}

	MongoDB = MongoClient
}
