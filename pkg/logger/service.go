package logger

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
)

// Collection here that we manage Logs db
var Collection *mongo.Collection

// DB poiner to mongoDB and can access from any where
var DB *mongo.Database

// Connect if db connected or not
var Connect bool

// OpenConnection function for open connection with mongodb once
func OpenConnection() {
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	err := client.Connect(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		Connect = false
		cancel()
	} else {
		Connect = true

		DB = client.Database("gateway")
		Collection = DB.Collection("log")
	}
	_ = cancel
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
