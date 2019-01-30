package logger

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"time"
)

var Collection *mongo.Collection
var DB *mongo.Database

func OpenConnection() {
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	DB = client.Database("gateway")

	Collection = DB.Collection("log")
}

// SetUserLog set the middleware logs here
func SetUserLog(log UserLog) {
	_, _ = Collection.InsertOne(context.Background(), log)
}

// SetSysLog set the middleware logs here
func SetSysLog(log SystemLog) {
	_, _ = Collection.InsertOne(context.Background(), log)
}


// ShowLogs query the database
func ShowLogs() {
	// bson.D{{}} return all records in database you can create a map here
	data1, _ := Collection.Find(context.Background(), bson.D{{}})

	for data1.Next(context.Background()) {
		raw, err := data1.DecodeBytes()
		if err != nil { log.Fatal(err) }
		fmt.Println(raw)
	}
}
