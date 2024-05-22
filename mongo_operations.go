package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveVideoChunksMongo(collection *mongo.Collection, chunkData []byte, i int) {
	hexData := fmt.Sprintf("%x", chunkData)
	name := fmt.Sprintf("video%d", i)
	documents := []interface{}{}
	for j := 1; j <= RowsPerQuery; j++ {
		doc := bson.D{
			{Key: "name", Value: name},
			{Key: "data", Value: hexData},
		}
		documents = append(documents, doc)
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadVideoChunksMongo(collection *mongo.Collection, videoName string) {
	ctx := context.TODO()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cursor.Close(ctx)
	documentsFound := cursor.RemainingBatchLength()
	if documentsFound == 0 {
		log.Printf("No documents found for video name %s", videoName)
		return
	}
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
			return
		}
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func UpdateVideoDataMongo(collection *mongo.Collection, chunkData []byte, i int) {
	hexData := fmt.Sprintf("%x", chunkData)
	name := fmt.Sprintf("video%d", i)
	filter := bson.D{{Key: "name", Value: name}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "data", Value: hexData}}}}
	opts := options.Update().SetUpsert(true) // Create a new document if one doesn't exist
	result, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount == 0 {
		fmt.Println("No matching document found, a new document will be inserted.")
	}
}

func DropVideoDataMongo(collection *mongo.Collection, i int) {
	name := fmt.Sprintf("video%d", i)
	filter := bson.D{{Key: "name", Value: name}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	if result.DeletedCount == 0 {
		fmt.Println("No matching document found to delete.")
	} else {
		fmt.Printf("Successfully deleted document with name %s\n", name)
	}
}

func FindVideoDataMongo(collection *mongo.Collection, i int) {
	name := fmt.Sprintf("video%d", i)
	filter := bson.D{{Key: "name", Value: name}}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No matching document found.")
		} else {
			log.Fatal(err)
		}
		return
	}
}

/*
func SaveVideoChunkMongo(collection *mongo.Collection, chunkData []byte) {
	// Insert a single document
	hexData := fmt.Sprintf("%x", chunkData)
	fmt.Println(hexData)
	//documents := []interface{}{
	//	bson.D{
	//		{Key: "name", Value: "video"},
	//		{Key: "data", Value: "hexData"},
	//	},
	//	}
	_, err := collection.InsertOne(context.TODO(), bson.D{
		{Key: "name", Value: "video"},
		{Key: "data", Value: hexData},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted chunk into database.")
}

func saveVideoChunkMongounused(collection *mongo.Collection, chunkData []byte) {
	// Insert a single document
	hexData := fmt.Sprintf("%x", chunkData)
	fmt.Println(hexData)
	//documents := []interface{}{
	//	bson.D{
	//		{Key: "name", Value: "video"},
	//		{Key: "data", Value: "hexData"},
	//	},
	//	}
	_, err := collection.InsertOne(context.TODO(), bson.D{
		{Key: "name", Value: "video"},
		{Key: "data", Value: hexData},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted chunk into database.")
}
*/
