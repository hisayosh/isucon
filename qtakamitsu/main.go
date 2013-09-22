package main

import (
  "bufio"
  "fmt"
  "os"
  "io"
  "net/http"

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

var db_conn *sql.DB

func rootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "root:\n")
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/artist/"):]
    fmt.Fprintf(w, "artist:\n")
    fmt.Fprintf(w, "  id: %s\n", id)
}

func ticketHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/ticket/"):]
    fmt.Fprintf(w, "artist:\n")
    fmt.Fprintf(w, "  id: %s\n", id)
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
    variation_id := r.PostFormValue("variation_id")
    member_id := r.PostFormValue("member_id")

    fmt.Fprintf(w, "buy:\n")
    fmt.Fprintf(w, "  variation_id: %s\n", variation_id)
    fmt.Fprintf(w, "  member_id: %s\n", member_id)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
    sql_file := "../config/database/initial_data.sql"

    file, err := os.Open(sql_file)
    defer file.Close()
    if err != nil {
        panic(err)
    }

    bufReader := bufio.NewReader(file)
    for i := 0; ; i++ {
        s, err := bufReader.ReadString('\n')
        if err == io.EOF {
            break
        }

        if s == "\n" {
            continue
        }
        //fmt.Fprintf(w, "[%s]\n", s)

        _, err = db_conn.Exec(s)
        if err != nil {
            fmt.Fprintf(w, "Error: admin exec failed.\n")
        }
    }
}

func adminCsvHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "adminCsv:\n")
}

func main() {
    db_string := fmt.Sprintf("%s:%s@/%s?charset=utf8", conf.Username, conf.Password, conf.Dbname)
    tmp_conn, err := sql.Open("mysql", db_string)
    if err != nil {
        fmt.Print("error: open failed.\n");
        return
    }
    db_conn = tmp_conn

    http.HandleFunc("/artist/", artistHandler);
    http.HandleFunc("/ticket/", ticketHandler);
    http.HandleFunc("/buy/", buyHandler);
    http.HandleFunc("/admin/", adminHandler);
    http.HandleFunc("/admin/order.csv", adminCsvHandler);
    http.HandleFunc("/", rootHandler);

    http.ListenAndServe("localhost:9999", nil)

    //defer db_conn.Close()
}

