package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	Collection *mongo.Collection
}

func (mo *MongoDB) SaveVideoChunk(chunkData []byte, i, connID, rowsPerQuery int, primaryKey bool) {
	hexData := fmt.Sprintf("%x", chunkData)
	name := fmt.Sprintf("v%d-%d", connID, i)
	documents := []interface{}{}
	keyName := "_id"
	if !primaryKey {
		keyName = "name"
	}
	for j := 1; j <= rowsPerQuery; j++ {
		doc := bson.D{
			{Key: keyName, Value: name},
			{Key: "data", Value: hexData},
		}
		documents = append(documents, doc)
	}
	_, err := mo.Collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Fatal(err)
	}
}

func (mo *MongoDB) ReadAllVideoChunks() {
	ctx := context.TODO()
	cursor, err := mo.Collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cursor.Close(ctx)
	documentsFound := cursor.RemainingBatchLength()
	if documentsFound == 0 {
		log.Printf("No documents found")
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

func (mo *MongoDB) UpdateVideoChunk(chunkData []byte, i, connID int, primaryKey bool) {
	hexData := fmt.Sprintf("%x", chunkData)
	name := fmt.Sprintf("v%d-%d", connID, i)
	keyName := "_id"
	if !primaryKey {
		keyName = "name"
	}
	filter := bson.D{{Key: keyName, Value: name}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "data", Value: hexData}}}}
	result, err := mo.Collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount == 0 {
		fmt.Println("No matching document found, a new document will be inserted.")
	}
}

func (mo *MongoDB) DropVideoChunk(i, connID int, primaryKey bool) {
	name := fmt.Sprintf("v%d-%d", connID, i)
	keyName := "_id"
	if !primaryKey {
		keyName = "name"
	}
	filter := bson.D{{Key: keyName, Value: name}}
	result, err := mo.Collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	if result.DeletedCount == 0 {
		fmt.Println("No matching document found to delete.")
	}
}

func (mo *MongoDB) ReadVideoChunk(i, connID int, primaryKey bool) {
	ctx := context.TODO()
	name := fmt.Sprintf("v%d-%d", connID, i)
	keyName := "_id"
	if !primaryKey {
		keyName = "name"
	}
	filter := bson.D{{Key: keyName, Value: name}}
	cursor, err := mo.Collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cursor.Close(ctx)
	documentsFound := cursor.RemainingBatchLength()
	if documentsFound == 0 {
		log.Printf("No documents found")
		return
	}
	for cursor.Next(ctx) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			log.Println(err)
		}
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func (mo *MongoDB) ReadVideoChunkold(i, connID int, primaryKey bool) {
	name := fmt.Sprintf("v%d-%d", connID, i)
	keyName := "_id"
	if !primaryKey {
		keyName = "name"
	}
	filter := bson.D{{Key: keyName, Value: name}}
	var result bson.M
	err := mo.Collection.FindOne(context.TODO(), filter).Decode(&result)
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
