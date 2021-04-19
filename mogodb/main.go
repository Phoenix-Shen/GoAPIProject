package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri string = "mongodb+srv://root:asdf@ssk.3hxej.mongodb.net/GOAPIPROJDB?retryWrites=true&w=majority&authSource=admin"

type User struct {
	Name string
	Age  int
}

func connectToMongoDB(uri string, name string, timeout time.Duration, num uint64, collectionName string) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	o := options.Client().ApplyURI(uri)
	o.SetMaxPoolSize(num)
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, err
	}

	return client.Database(name).Collection(collectionName), nil
}

func Create(collection *mongo.Collection, obj []interface{}) []interface{} {
	insertResult, err := collection.InsertMany(context.TODO(), obj)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return insertResult.InsertedIDs
}

func Read(collection *mongo.Collection /*, obj []interface{}*/) []interface{} {

	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var results []interface{}
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		var elem interface{}
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// 完成后关闭游标
	cur.Close(context.TODO())
	return results
}

func Update(collection *mongo.Collection, obj *interface{}) {
	log.Fatal("not implementd Error")
}

func main() {
	/*ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://root:asdf@ssk.3hxej.mongodb.net/GOAPIPROJDB?retryWrites=true&w=majority&authSource=admin",
	))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("connected to MongoDB")*/
	conn, _ := connectToMongoDB(uri, "sample_mflix", 10*time.Second, 5, "comments")
	print(Read(conn))

}

//
