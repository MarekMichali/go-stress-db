package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type MysqlDB struct {
	Db *sql.DB
}

func (my *MysqlDB) SaveVideoChunksExp(chunkData []byte, i, rowsPerQuery int) {
	hexData := fmt.Sprintf("%x", chunkData)
	name := fmt.Sprintf("video%d", i)
	var sb strings.Builder
	sb.WriteString("INSERT INTO videos (name, data) VALUES ")
	for j := 1; j <= rowsPerQuery; j++ {
		sb.WriteString(fmt.Sprintf("('%s', X'%x'),", name, hexData))
	}
	query := strings.TrimRight(sb.String(), ",")
	_, err := my.Db.Exec(query)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}

func (my *MysqlDB) ReadVideoChunk() {
	_, err := my.Db.Exec("SELECT * FROM videos")
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No rows found ")
		} else {
			log.Fatal(err)
		}
		return
	}
}

func (my *MysqlDB) UpdateVideoData(chunkData []byte, i int) {
	hexData := fmt.Sprintf("%x", chunkData)
	name := fmt.Sprintf("video%d", i)
	query := fmt.Sprintf("UPDATE videos SET data=X'%x' where name='%s'", hexData, name)
	_, err := my.Db.Exec(query)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}

func (my *MysqlDB) DropVideoData(i int) {
	name := fmt.Sprintf("video%d", i)
	query := fmt.Sprintf("DELETE FROM videos where name='%s'", name)
	_, err := my.Db.Exec(query)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}

func (my *MysqlDB) FindVideoChunks(i int) {
	name := fmt.Sprintf("video%d", i)
	query := fmt.Sprintf("SELECT * FROM videos where name='%s'", name)
	_, err := my.Db.Exec(query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No rows found ")
		} else {
			log.Fatal(err)
		}
		return
	}
}

/*
func DeleteVideoChunk(db *sql.DB) {
	_, err := db.Exec("DELETE FROM videos")
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	log.Println("Deleted video chunk from the database.")
}

func updateVideoChunk(db *sql.DB, videoName string, chunkData []byte) {
	_, err := db.Exec("UPDATE videos SET data = ? WHERE name = ?", chunkData, videoName)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	log.Println("Updated video chunk in the database.")
}
*/
