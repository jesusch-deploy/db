package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

type DbSingleton struct {
	connection *sql.DB
}

var (
	dbInstance *DbSingleton
	once       sync.Once

// counts     = 0
)

const dbTimeout = time.Second * 3

func GetInstance(dsn string) *DbSingleton {
	once.Do(func() {
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			log.Printf("Base de datos no esta lista %v", err)
		}
		dbInstance = &DbSingleton{
			connection: db,
		}
	})
	return dbInstance
}

func (db *DbSingleton) Close() {
	if db.connection != nil {
		err := db.connection.Close()
		if err != nil {
			log.Printf("Error al cerrar la base de datos %v", err)
		}
	}
}

func (db *DbSingleton) GetConnection() *sql.DB {
	return db.connection
}

/*
func OpenDatabase(dsn string) *sql.DB {
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Base de datos no esta lista...")
			counts++
		} else {
			log.Println("Base de datos conectada")
			dbInstance = connection
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
*/

func Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := dbInstance.connection.Exec(query, args...)
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := dbInstance.connection.Query(query, args)
	if err != nil {
		fmt.Println(err)
	}
	return rows, err
}

func QueryRowContext(query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var newId int
	err := dbInstance.connection.QueryRowContext(ctx, query, args).Scan(&newId)
	if err != nil {
		log.Println("Error en el registro:", err)
		return err
	}
	return nil
}

/*
func Close() {
	dbInstance.Close()
}
*/
