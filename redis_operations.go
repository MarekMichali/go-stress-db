package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Client *redis.Client
}

func (rdb *RedisDB) SaveVideoChunk(chunkData []byte, i, connID int) {
	hexData := fmt.Sprintf("%x", chunkData)
	ctx := context.Background()
	key := fmt.Sprintf("v%d-%d", connID, i)
	err := rdb.Client.Set(ctx, key, hexData, 0).Err()
	if err != nil {
		log.Fatalf("Unable to save the chunk to Redis. %v", err)
	}
}

func (rdb *RedisDB) ReadAllVideoChunks() {
	ctx := context.TODO()
	keys, err := rdb.Client.Keys(ctx, "v*").Result()
	if err != nil {
		log.Fatalf("Unable to retrieve keys from Redis. %v", err)
		return
	}
	if len(keys) == 0 {
		log.Printf("No documents found")
		return
	}
	for _, key := range keys {
		_, err := rdb.Client.Get(ctx, key).Bytes()
		if err != nil {
			log.Fatalf("Unable to retrieve chunk data from Redis. %v", err)
			return
		}
	}
}

func (rdb *RedisDB) UpdateVideoChunk(chunkData []byte, i, connID int) {
	hexData := fmt.Sprintf("%x", chunkData)
	key := fmt.Sprintf("v%d-%d", connID, i)
	err := rdb.Client.Set(context.TODO(), key, hexData, 0).Err()
	if err != nil {
		log.Fatalf("Unable to update the chunk in Redis. %v", err)
	}
}

func (rdb *RedisDB) DropVideoChunk(i, connID int) {
	key := fmt.Sprintf("v%d-%d", connID, i)
	err := rdb.Client.Del(context.TODO(), key).Err()
	if err != nil {
		log.Fatalf("Unable to delete the video chunk from Redis. %v", err)
	}
}

func (rdb *RedisDB) ReadVideoChunk(i, connID int) {
	ctx := context.TODO()
	key := fmt.Sprintf("v%d-%d", connID, i)
	_, err := rdb.Client.Get(ctx, key).Bytes()
	if err != nil {
		log.Fatalf("Unable to retrieve chunk data from Redis. %v", err)
		return
	}
}
