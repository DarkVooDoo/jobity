package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var database *sql.DB

func InitDB()error{
    connString := fmt.Sprintf("postgres://darkvoodoo:%v@db:5432/%v?sslmode=disable", os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
    db, err := sql.Open("postgres", connString)
    if err != nil{
        return errors.New("error conn")
    }
    db.SetMaxIdleConns(5)
    db.SetConnMaxIdleTime(time.Second*10)
    db.SetMaxOpenConns(15)
    db.SetConnMaxLifetime(time.Second * 10)
    database = db 
    return nil
}

func GetDBConn()(*sql.Conn, error){
    if database == nil{
        return nil, errors.New("database pointer is nil")
    } 
    return database.Conn(context.Background())
}

func GetDBPoolConn() *sql.Conn{
    cxt := context.Background()
    conn, err := database.Conn(cxt)
    if err != nil{
        log.Println("error in the conn")
        return conn
    }
    log.Printf("Conn in use: %v", database.Stats().InUse)
    return conn
}
