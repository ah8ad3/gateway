package logger

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// Collection here that we manage Logs db
var Collection *mongo.Collection

// DB poiner to mongoDB and can access from any where
var DB *mongo.Database

// Connect if db connected or not
var Connect bool

// OpenConnection function for open connection with mongodb once
func OpenConnection() {
	// client, err := mongo.NewClientWithOptions("mongodb://localhost:27017")
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		fmt.Println(err.Error())
		Connect = false
	} else {
		Connect = true
		err = client.Connect(context.Background())

		DB = client.Database("gateway")

		Collection = DB.Collection("log")
	}
}

// SetUserLog set the middleware logs here
func SetUserLog(log UserLog) {
	if Connect {
		_, _ = Collection.InsertOne(context.Background(), log)
	}
}

// SetSysLog set the middleware logs here
func SetSysLog(log SystemLog) {
	if Connect {
		_, _ = Collection.InsertOne(context.Background(), log)
	}
}

// ShowLogs query the database
// func ShowLogs() {
// 	// bson.D{{}} return all records in database you can create a map here
// 	data1, _ := Collection.Find(context.Background(), bson.D{{}})

// 	for data1.Next(context.Background()) {
// 		raw, err := data1.DecodeBytes()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(raw)
// 	}
// }
