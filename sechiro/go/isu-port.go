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

    db, err := sql.Open("mysql", "root@/isucon2")
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()
    bufReader := bufio.NewReader(file)
    for i := 0; ; i++ {
        s, err := bufReader.ReadString('\n')
        if err == io.EOF {
            break
        }
        fmt.Fprintf(w, "%s\n", s)
        db.Exec(s)
    }
    // Prepare statement for reading data
    stmtOut, err := db.Prepare("SELECT id FROM stock WHERE id = ?")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtOut.Close()

    var squareNum int
    // Query the square-number of 13
    err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    fmt.Printf("The square number of 13 is: %d", squareNum)
}

func main() {
    http.HandleFunc("/artist/", artistHandler);
    http.HandleFunc("/ticket/", ticketHandler);
    http.HandleFunc("/buy/", buyHandler);
    http.HandleFunc("/admin/", adminHandler);
    http.HandleFunc("/", rootHandler);

    http.ListenAndServe(":8080", nil)
}
