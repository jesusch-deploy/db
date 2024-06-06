package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

var db *sql.DB
var counts = 0

const dbTimeout = time.Second * 3

func OpenDatabase(dsn string) *sql.DB {
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Base de datos no esta lista...")
			counts++
		} else {
			log.Println("Base de datos conectada")
			db = connection
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Intentado en dos segundos...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args)
	if err != nil {
		fmt.Println(err)
	}
	return rows, err
}

func QueryRowContext(query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var newId int
	err := db.QueryRowContext(ctx, query, args).Scan(&newId)
	if err != nil {
		log.Println("Error en el registro:", err)
		return err
	}
	return nil
}

func Close() {
	db.Close()
}
