package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"homework/database"
	"log"
	"time"
)

type User struct {
	UserId         string
	Message        string
	Name           string
	PictureURL     string
	ProfileMessage string
}

type UserInfo struct {
	UserId string
	Name   string
}

func InsertDataToDBCollection(db string, collection string, data interface{}) error {

	dbcollection := database.MongoDB.Database(db).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := dbcollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func GetAllUserInfo() (users []UserInfo, err error) {

	tmp := make(map[string]bool)

	dbcollection := database.MongoDB.Database("db").Collection("user_info")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{"userid", 1}})
	cur, err := dbcollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		var user UserInfo
		err := cur.Decode(&user)
		if err != nil {
			return users, err
		}

		if _, ok := tmp[user.UserId]; !ok {
			users = append(users, user)
			tmp[user.UserId] = true
		}
	}

	return
}

func GetUserInfo(userid string) (users []User, err error) {
	dbcollection := database.MongoDB.Database("db").Collection("user_info")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := dbcollection.Find(ctx, bson.M{"userid": userid})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return
}
