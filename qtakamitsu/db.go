package main

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type Config struct {
    Host string
    Port int
    Username string
    Password string
    Dbname string
}

var conf = Config {
    Host: "127.0.0.1",
    Port: 3306,
    Username: "isucon2app",
    Password: "isunageruna",
    Dbname: "isucon2" }

var db_conn sql.DB

func main() {
    db_string := fmt.Sprintf("%s:%s@/%s", conf.Username, conf.Password, conf.Dbname)
    db_conn, err := sql.Open("mysql", db_string)
    if err != nil {
        fmt.Print("error: open failed.\n");
        return
    }
    defer db_conn.Close()

    fmt.Print("success open.\n");
}

