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

	"flag"
)

var (
	BufferSize      int
	NoOfConnections int
	TickInterval    time.Duration
	RowsPerQuery    int
	SelectDB        string
	OpType          string
	MaxIterations   int
	mysqlConnStr    string
	mariadbConnStr  string
	mongoConnStr    string
	redisConnStr    string
	videoName       string
	primaryKey      bool
)

func flags() {
	flag.IntVar(&BufferSize, "buffer", 8000, "Buffer size")
	flag.IntVar(&NoOfConnections, "conn", 1, "Number of connections")
	flag.DurationVar(&TickInterval, "intvl", 1*time.Millisecond, "Tick interval between queries")
	flag.IntVar(&RowsPerQuery, "rows", 1, "Rows per query")
	flag.StringVar(&SelectDB, "db", "mysql", "Select database (mysql, mariadb, mongodb, redis)")
	flag.StringVar(&OpType, "op", "insert", "Operation type (insert, select, update, delete)")
	flag.IntVar(&MaxIterations, "it", 50000, "Max iterations")
	flag.StringVar(&mysqlConnStr, "mysqlConnStr", "root:123456@tcp(127.0.0.1:3333)/videos", "MySQL connection string")
	flag.StringVar(&mariadbConnStr, "mariadbConnStr", "root:123456@tcp(127.0.0.1:3306)/videos", "MariaDB connection string")
	flag.StringVar(&mongoConnStr, "mongoConnStr", "mongodb://pmm:pmm@localhost:27017/?serverSelectionTimeoutMS=30000", "MongoDB connection string")
	flag.StringVar(&redisConnStr, "redisConnStr", "redis://localhost:6379", "Redis connection string")
	flag.StringVar(&videoName, "video", "bigSample.mp4", "Video name")
	flag.BoolVar(&primaryKey, "pk", true, "Use primary key mode")
}

func main() {
	flags()
	flag.Parse()
	var wg sync.WaitGroup
	fmt.Printf("Starting, Timestamp: %s\n", time.Now().Format(time.StampMilli))
	fmt.Printf("Database selected: %s, Operation type: %s, Number of connections: %d, Iterations: %d, Rows per query: %d, Buffer size: %d, Query interval: %s, Primary key mode: %v\n", SelectDB, OpType, NoOfConnections, MaxIterations, RowsPerQuery, BufferSize, TickInterval, primaryKey)
	if SelectDB == "mysql" || SelectDB == "mariadb" {
		if SelectDB == "mariadb" {
			mysqlConnStr = mariadbConnStr
		}
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
				if OpType == "insert" {
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
					/*		} else if OpType == 2 {
							for range ticker.C {
								if i >= MaxIterations {
									fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
									break
								}
								mysqlDB.ReadAllVideoChunks()
								i++
							}*/
				} else if OpType == "update" {
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
				} else if OpType == "delete" {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mysqlDB.DropVideoChunk(i, connID)
						i++
					}
				} else if OpType == "select" {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mysqlDB.ReadVideoChunk(i, connID)
						i++
					}
				} else {
					log.Fatal("Invalid operation type selected")
				}
			}(connID)
		}
	} else if SelectDB == "mongodb" {
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
				if OpType == "insert" {
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
						mongoDB.SaveVideoChunk(buffer[:bytesRead], i, connID, RowsPerQuery, primaryKey)
						i++
					}
					/*} else if OpType == 2 {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mongoDB.ReadAllVideoChunks()
						i++
					}*/
				} else if OpType == "update" {
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
						mongoDB.UpdateVideoChunk(buffer[:bytesRead], i, connID, primaryKey)
						i++
					}
				} else if OpType == "delete" {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mongoDB.DropVideoChunk(i, connID, primaryKey)
						i++
					}
				} else if OpType == "select" {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						mongoDB.ReadVideoChunk(i, connID, primaryKey)
						i++
					}
				} else {
					log.Fatal("Invalid operation type selected")
				}
			}(connID)
		}
	} else if SelectDB == "redis" {
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
				if OpType == "insert" {
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
				} else if OpType == "select" {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						RedisDB.ReadVideoChunk(i, connID)
						i++
					}
				} else if OpType == "update" {
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
				} else if OpType == "delete" {
					for range ticker.C {
						if i >= MaxIterations {
							fmt.Printf("Max Iterations reached, Timestamp: %s\n", time.Now().Format(time.StampMilli))
							break
						}
						RedisDB.DropVideoChunk(i, connID)
						i++
					}
				} else {
					log.Fatal("Invalid operation type selected")
				}
			}(connID)

		}
	} else {
		log.Fatal("Invalid database selected")
	}
	wg.Wait()
}
