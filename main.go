package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BufferSize      = 8000 //26350 vor varbinary, 8000 for varchar // 61k iteration for 8000 for mysql
	NoOfConnections = 1
	TickInterval    = 1 * time.Millisecond
	RowsPerQuery    = 1     // works for inserts only
	SelectDB        = 3     // 1 for mysql, 2 for mongo, 3 for redis
	OpType          = 4     // 1 for insert, 2 for select all 3 for update 4 for delete, 5 for select 1 row
	MaxIterations   = 50000 // 50000 max //2k for 25 conn
	mysqlConnStr    = "root:123456@tcp(127.0.0.1:3306)/videos"
	mongoConnStr    = "mongodb://pmm:pmm@localhost:27017/?serverSelectionTimeoutMS=3000000"
	redisConnStr    = "redis://localhost:6379"
	videoName       = "bigSample.mp4"
)

func main() {
	var wg sync.WaitGroup
	fmt.Printf("Starting, Timestamp: %s\n", time.Now().Format(time.StampMilli))
	fmt.Printf("Operation type: %d, max iterations: %d, buffer size: %d, rows per query: %d, no of connections: %d\n", OpType, MaxIterations, BufferSize, RowsPerQuery, NoOfConnections)
	if SelectDB == 1 {
		for connID := 0; connID < NoOfConnections; connID++ {
			wg.Add(1)
			go func(connID int) {
				defer wg.Done()

				file, err := os.Open(videoName)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				db, err := sql.Open("mysql", mysqlConnStr)
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()

				mysqlDB := MysqlDB{
					Db: db,
				}
				buffer := make([]byte, BufferSize)
				ticker := time.NewTicker(TickInterval)
				defer ticker.Stop()

				i := 0
				if OpType == 1 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						bytesRead, err := file.Read(buffer)
						if err != nil {
							if err != io.EOF {
								log.Fatal(err)
							}
							return
						}
						//	bufCopy := make([]byte, bytesRead)
						//	copy(bufCopy, buffer[:bytesRead])
						//	go mysqlDB.SaveVideoChunk(bufCopy, j, RowsPerQuery)
						mysqlDB.SaveVideoChunk(buffer[:bytesRead], i, connID, RowsPerQuery)
						i++
					}
				} else if OpType == 2 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mysqlDB.ReadAllVideoChunks()
						i++
					}
				} else if OpType == 3 {
					file.Read(buffer)
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						bytesRead, err := file.Read(buffer)
						if err != nil {
							if err != io.EOF {
								log.Fatal(err)
							}
							return
						}
						//	bufCopy := make([]byte, bytesRead)
						//	copy(bufCopy, buffer[:bytesRead])
						//	go mysqlDB.UpdateVideoChunk(bufCopy, j)
						mysqlDB.UpdateVideoChunk(buffer[:bytesRead], i, connID)
						i++
					}
				} else if OpType == 4 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mysqlDB.DropVideoChunk(i, connID)
						i++
					}
				} else if OpType == 5 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mysqlDB.ReadVideoChunk(i, connID)
						i++
					}
				}
			}(connID)
		}
	} else if SelectDB == 2 {
		for connID := 0; connID < NoOfConnections; connID++ {
			wg.Add(1)
			go func(connID int) {
				defer wg.Done()
				file, err := os.Open(videoName)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				clientOptions := options.Client().ApplyURI(mongoConnStr)
				client, err := mongo.Connect(context.TODO(), clientOptions)
				if err != nil {
					log.Fatal(err)
				}
				defer func() {
					if err = client.Disconnect(context.TODO()); err != nil {
						panic(err)
					}
				}()

				collection := client.Database("test").Collection("videos")
				mongoDB := MongoDB{
					Collection: collection,
				}
				buffer := make([]byte, BufferSize)
				ticker := time.NewTicker(TickInterval)
				defer ticker.Stop()

				i := 0
				if OpType == 1 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						bytesRead, err := file.Read(buffer)
						if err != nil {
							if err != io.EOF {
								log.Fatal(err)
							}
							return
						}
						//	bufCopy := make([]byte, bytesRead)
						//	copy(bufCopy, buffer[:bytesRead])
						//	go mongoDB.SaveVideoChunk(bufCopy, j, RowsPerQuery)
						mongoDB.SaveVideoChunk(buffer[:bytesRead], i, connID, RowsPerQuery)
						i++
					}
				} else if OpType == 2 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mongoDB.ReadAllVideoChunks()
						i++
					}
				} else if OpType == 3 {
					file.Read(buffer)
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						bytesRead, err := file.Read(buffer)
						if err != nil {
							if err != io.EOF {
								log.Fatal(err)
							}
							return
						}
						//	bufCopy := make([]byte, bytesRead)
						//	copy(bufCopy, buffer[:bytesRead])
						//	go mongoDB.UpdateVideoChunk(bufCopy, j)
						mongoDB.UpdateVideoChunk(buffer[:bytesRead], i, connID)
						i++
					}
				} else if OpType == 4 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mongoDB.DropVideoChunk(i, connID)
						i++
					}
				} else if OpType == 5 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mongoDB.ReadVideoChunk(i, connID)
						i++
					}
				}
			}(connID)
		}
	} else {
		for connID := 0; connID < NoOfConnections; connID++ {
			wg.Add(1)
			go func(connID int) {
				defer wg.Done()
				file, err := os.Open(videoName)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				opts, err := redis.ParseURL(redisConnStr)
				if err != nil {
					panic(err)
				}

				client := redis.NewClient(opts)
				RedisDB := RedisDB{
					Client: client,
				}

				buffer := make([]byte, BufferSize)
				ticker := time.NewTicker(TickInterval)
				defer ticker.Stop()

				i := 0
				if OpType == 1 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						bytesRead, err := file.Read(buffer)
						if err != nil {
							if err != io.EOF {
								log.Fatal(err)
							}
							return
						}
						RedisDB.SaveVideoChunk(buffer[:bytesRead], i, connID)
						i++
					}
				} else if OpType == 5 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						RedisDB.ReadVideoChunk(i, connID)
						i++
					}
				} else if OpType == 3 {
					file.Read(buffer)
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						bytesRead, err := file.Read(buffer)
						if err != nil {
							if err != io.EOF {
								log.Fatal(err)
							}
							return
						}

						RedisDB.UpdateVideoChunk(buffer[:bytesRead], i, connID)
						i++
					}
				} else if OpType == 4 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						RedisDB.DropVideoChunk(i, connID)
						i++
					}
				}
			}(connID)

		}
	}
	wg.Wait()
}
