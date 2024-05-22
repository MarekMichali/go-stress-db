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
	RowsPerQuery    = 1     // works for inserts only
	SelectDB        = 1     // 1 for mysql, 2 for mongo
	OpType          = 1     // 1 for insert, 2 for select all 3 for update 4 for delete, 5 for select 1 row
	MaxIterations   = 50000 // 50000 max
	mysqlConnStr    = "root:123456@tcp(127.0.0.1:3306)/videos"
	mongoConnStr    = "mongodb://root:123456@localhost:27017"
	videoName       = "bigSample.mp4"
)

func main() {
	var wg sync.WaitGroup
	fmt.Printf("Starting, Timestamp: %s\n", time.Now().Format(time.StampMilli))

	if SelectDB == 1 {
		for i := 0; i < NoOfConnections; i++ {
			wg.Add(1)
			go func() {
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
						mysqlDB.SaveVideoChunk(buffer[:bytesRead], j, RowsPerQuery)
						j++
					}
				} else if OpType == 2 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mysqlDB.ReadAllVideoChunks()
						j++
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
						mysqlDB.UpdateVideoChunk(buffer[:bytesRead], j)
						j++
					}
				} else if OpType == 4 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mysqlDB.DropVideoChunk(j)
						j++
					}
				} else if OpType == 5 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mysqlDB.ReadVideoChunk(j)
						j++
					}
				}
			}()
		}
	} else {
		for i := 0; i < NoOfConnections; i++ {
			wg.Add(1)
			go func() {
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
						mongoDB.SaveVideoChunk(buffer[:bytesRead], j, RowsPerQuery)
						j++
					}
				} else if OpType == 2 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mongoDB.ReadAllVideoChunks()
						j++
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
						mongoDB.UpdateVideoChunk(buffer[:bytesRead], j)
						j++
					}
				} else if OpType == 4 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mongoDB.DropVideoChunk(j)
						j++
					}
				} else if OpType == 5 {
					for range ticker.C {
						if j >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mongoDB.ReadVideoChunk(j)
						j++
					}
				}
			}()
		}
	}
	wg.Wait()
}
