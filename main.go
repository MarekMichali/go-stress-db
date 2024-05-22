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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BufferSize      = 8000 //26350 vor varbinary, 8000 for varchar // 61k iteration for 8000 for mysql
	NoOfConnections = 1
	TickInterval    = 1 * time.Millisecond
	RowsPerQuery    = 1
	SelectDB        = 1     // 1 for mysql, 2 for mongo
	OpType          = 1     // 1 for insert, 2 for select all 3 for update 4 for delete, 5 for select 1 row
	MaxIterations   = 50000 // 50000 max
)

func main() {
	var wg sync.WaitGroup
	fmt.Printf("Starting, Timestamp: %s\n", time.Now().Format(time.StampMilli))

	if SelectDB == 1 {
		for i := 0; i < NoOfConnections; i++ {
			x := i
			wg.Add(1)
			go func() {
				defer wg.Done()

				file, err := os.Open("bigSample.mp4")
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/videos")
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()

				buffer := make([]byte, BufferSize)
				ticker := time.NewTicker(TickInterval)
				defer ticker.Stop()
				j := 0
				if OpType == 1 {

					for range ticker.C {
						if j >= MaxIterations {
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
						//ReadVideoChunk(db, "video")
						//log.Println("readed")
						//return
						//DeleteVideoChunk(db)
						//return
						x = j
						SaveVideoChunksExp(db, buffer[:bytesRead], x)
						j++
						//return
					}
				} else if OpType == 2 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}

						//readVideoChunks(db, "video")
						//x = j
						ReadVideoChunk(db)
						j++
						//return
					}
				} else if OpType == 3 {
					file.Read(buffer)
					for range ticker.C {
						if j >= MaxIterations {
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
						x = j
						UpdateVideoData(db, buffer[:bytesRead], x)
						j++
						//return
					}
				} else if OpType == 4 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						x = j
						DropVideoData(db, x)
						j++
						//return
					}
				} else if OpType == 5 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}

						//readVideoChunks(db, "video")
						x = j
						FindVideoChunks(db, x)
						j++
						//return
					}
				}
			}()
		}
	} else {
		for i := 0; i < NoOfConnections; i++ {
			x := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				file, err := os.Open("bigSample.mp4")
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				clientOptions := options.Client().ApplyURI("mongodb://root:123456@localhost:27017")
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
				buffer := make([]byte, BufferSize)
				ticker := time.NewTicker(TickInterval)
				defer ticker.Stop()
				j := 0
				if OpType == 1 {
					for range ticker.C {
						if j >= MaxIterations {
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
						x = j
						SaveVideoChunksMongo(collection, buffer[:bytesRead], x)
						j++
						//return
					}
				} else if OpType == 2 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}

						ReadVideoChunksMongo(collection, "video")
						//x = j
						//FindVideoDataMongo(collection, x)
						j++
						//return
					}
				} else if OpType == 3 {
					file.Read(buffer)
					for range ticker.C {
						if j >= MaxIterations {
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
						x = j
						UpdateVideoDataMongo(collection, buffer[:bytesRead], x)
						j++
						//return
					}
				} else if OpType == 4 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						x = j
						DropVideoDataMongo(collection, x)
						j++
						//return
					}
				} else if OpType == 5 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}

						//ReadVideoChunksMongo(collection, "video")
						x = j
						FindVideoDataMongo(collection, x)
						j++
						//return
					}
				}
			}()
		}
	}
	wg.Wait()
}
