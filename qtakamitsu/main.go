package main

import (
  "os"
  "bufio"
  "io"
  "fmt"
  "text/template"
  "net/http"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

////////////////////////////////////////////////////////////////////////////////
var db DbConn

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

type DbConn struct {
    db_conn *sql.DB
}

////////////////////////////////////////////////////////////////////////////////
func (p *DbConn) Init() {
    fmt.Print("DbConn: init\n");

    db_string := fmt.Sprintf("%s:%s@/%s", conf.Username, conf.Password, conf.Dbname)
    db_conn, err := sql.Open("mysql", db_string)
    if err != nil {
        return
    }

    p.db_conn = db_conn
}

func (p *DbConn) Close() {
    p.db_conn.Close()
}

func (p *DbConn) RecentSold() []map[string]string {
    fmt.Print("DbConn: RecentSold\n");

    sql := "" +
        "SELECT stock.seat_id, variation.name AS v_name, ticket.name AS t_name, artist.name AS a_name FROM stock" +
        " JOIN variation ON stock.variation_id = variation.id" +
        " JOIN ticket ON variation.ticket_id = ticket.id" +
        " JOIN artist ON ticket.artist_id = artist.id" +
        " WHERE order_id IS NOT NULL" +
        " ORDER BY order_id DESC LIMIT 10"

    fmt.Printf("  [%s]\n", sql);

    rows, err := p.db_conn.Query(sql)
    if err != nil {
        panic(err.Error())
    }
    defer rows.Close()

    return GetData(rows);
}

type IndexViewModel struct {
    RecentSold []map[string]string
    Artists []map[string]string
}

func (p *DbConn) Index() IndexViewModel {
    fmt.Print("DbConn: get /\n");
    sql := "SELECT * FROM artist ORDER BY id"
    fmt.Printf("  [%s]\n", sql);

    rows, err := p.db_conn.Query(sql)
    if err != nil {
        panic(err.Error())
    }
    //defer rows.Close()

    recent := db.RecentSold()
    artists := GetData(rows)

    model := IndexViewModel {
        RecentSold: recent,
        Artists: artists,
    }

    return model
}

func GetData(r *sql.Rows) []map[string]string {
    // Get column names
    columns, err := r.Columns()
    if err != nil {
        panic(err.Error())
    }

    // Make a slice for the values
    //values := make([]sql.RawBytes, len(columns))
    values := make([][]byte, len(columns))

    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    data := make([]map[string]string, 0)
    for r.Next() {
        // get RawBytes from data
        err = r.Scan(scanArgs...)
        if err != nil {
            panic(err.Error())
        }

        item := make(map[string]string)
        // Now do something with the data.
        // Here we just print each column as a string.
        var value string
        for i, col := range values {
            if col == nil {
                value = ""
            } else {
                value = string(col)
            }
            item[columns[i]] = value
        }

        data = append(data, item)
    }

    return data
}

func PrintData(data []map[string]string) {
    siz := len(data)
    if siz == 0 {
        fmt.Println("data empty.")
        return
    }

    for i := 0; i < siz; i++ {
        for k, v := range data[i] {
            fmt.Printf("%s=%s, ", k, v)
        }
        fmt.Print("\n")
    }
}

type ArtistViewModel struct {
    RecentSold []map[string]string

    Artist map[string]string
    Tickets []map[string]string
}

func (p *DbConn) GetArtists(artist_id string) ArtistViewModel {
    fmt.Printf("DbConn: get /artist/:artistid  (%s)\n", artist_id);

    sql_artist := "SELECT id, name FROM artist WHERE id = ? LIMIT 1"
    sql_ticket := "SELECT id, name FROM ticket WHERE artist_id = ? ORDER BY id"
    sql_count_tickets := "SELECT COUNT(*) as num FROM variation" +
        " INNER JOIN stock ON stock.variation_id = variation.id" +
        " WHERE variation.ticket_id = ? AND stock.order_id IS NULL"

    fmt.Printf("  [%s]\n", sql_artist);
    fmt.Printf("  [%s]\n", sql_ticket);
    fmt.Printf("  [%s]\n", sql_count_tickets);

    stmt, err := p.db_conn.Prepare(sql_artist)
    if err != nil {
        panic(err.Error())
    }
    rows, err := stmt.Query(artist_id)
    if err != nil {
        panic(err.Error())
    }
    defer rows.Close()
    artists := GetData(rows);


    stmt, err = p.db_conn.Prepare(sql_ticket)
    if err != nil {
        panic(err.Error())
    }
    rows, err = stmt.Query(artists[0]["id"])
    if err != nil {
        panic(err.Error())
    }
    defer rows.Close()
    tickets := GetData(rows);

    siz := len(tickets)
    for i := 0; i < siz; i++ {
        stmt, err = p.db_conn.Prepare(sql_count_tickets)
        if err != nil {
            panic(err.Error())
        }
        rows, err = stmt.Query(tickets[i]["id"])
        if err != nil {
            panic(err.Error())
        }
        count := GetData(rows);

        tickets[i]["count"] = count[0]["num"]
    }
    recent := db.RecentSold()

    model := ArtistViewModel {
        RecentSold: recent,
        Artist: artists[0],
        Tickets: tickets,
    }

    return model
}

type SeatInfo struct {
    Items []map[string]string
    Count string
    VariationId string
    VariationName string
}

type TicketViewModel struct {
    RecentSold []map[string]string

    Ticket map[string]string
    Variations []map[string]string

    Seat []SeatInfo
}

func (p *DbConn) GetTickets(ticket_id string) TicketViewModel {
    fmt.Printf("DbConn: get /ticket/:ticketid (%s)\n", ticket_id)

    sql_ticket := "SELECT t.*, a.name AS artist_name FROM ticket t INNER JOIN artist a ON t.artist_id = a.id WHERE t.id = ? LIMIT 1"
    sql_variations := "SELECT id, name FROM variation WHERE ticket_id = ? ORDER BY id"
    sql_variation_stock := "SELECT seat_id, order_id FROM stock WHERE variation_id = ?"
    sql_variation_vacancy := "SELECT COUNT(*) as num FROM stock WHERE variation_id = ? AND order_id IS NULL"

    fmt.Printf("  [%s]\n", sql_ticket);
    fmt.Printf("  [%s]\n", sql_variations);
    fmt.Printf("  [%s]\n", sql_variation_stock);
    fmt.Printf("  [%s]\n", sql_variation_vacancy);


    stmt, err := p.db_conn.Prepare(sql_ticket)
    if err != nil {
        panic(err.Error())
    }
    rows, err := stmt.Query(ticket_id)
    if err != nil {
        panic(err.Error())
    }
    tickets := GetData(rows);
    //PrintData(ticket)

    stmt, err = p.db_conn.Prepare(sql_variations)
    if err != nil {
        panic(err.Error())
    }
    rows, err = stmt.Query(tickets[0]["id"])
    if err != nil {
        panic(err.Error())
    }
    defer rows.Close()
    variations := GetData(rows);
    //PrintData(variations)

    seats := make([]SeatInfo, 0)

    siz := len(variations)
    for i := 0; i < siz; i++ {
        stmt, err = p.db_conn.Prepare(sql_variation_stock)
        if err != nil {
            panic(err.Error())
        }
        rows, err = stmt.Query(variations[i]["id"])
        if err != nil {
            panic(err.Error())
        }
        stack := GetData(rows);
        //PrintData(stack)

        siz := len(stack)
        for i := 0; i < siz; i++ {
            if stack[i]["order_id"] == "" {
                stack[i]["order_id"] = "available"
            } else {
                stack[i]["order_id"] = "unavailable"
            } 
        }

        seat := SeatInfo {
            Items: stack,
            VariationId: variations[i]["id"],
            VariationName: variations[i]["name"],
        }

        stmt, err = p.db_conn.Prepare(sql_variation_vacancy)
        if err != nil {
            panic(err.Error())
        }
        rows, err = stmt.Query(variations[i]["id"])
        if err != nil {
            panic(err.Error())
        }
        count := GetData(rows);
        //PrintData(num)

        seat.Count = count[0]["num"]

        seats = append(seats, seat)
    }
    recent := db.RecentSold()

    model := TicketViewModel {
        RecentSold: recent,
        Ticket: tickets[0],
        Variations: variations,
        Seat: seats,
    }

    return model
}

type BuyViewModel struct {
    RecentSold []map[string]string
    SeatId string
    MemberId string
    Result bool
}

func (p *DbConn) Buy(variation_id string, member_id string) BuyViewModel {
    fmt.Printf("DbConn: post /buy (%s:%s)\n", variation_id, member_id)

    sql_add := "INSERT INTO order_request (member_id) VALUES (?)"
    sql_mod := "UPDATE stock SET order_id = ? WHERE variation_id = ? AND order_id IS NULL ORDER BY RAND() LIMIT 0"
    sql_get := "SELECT seat_id FROM stock WHERE order_id = ? LIMIT 1', $order_id"

    fmt.Printf("  [%s]\n", sql_add);
    fmt.Printf("  [%s]\n", sql_mod);
    fmt.Printf("  [%s]\n", sql_get);

    recent := db.RecentSold()

    tx, err := p.db_conn.Begin()

    stmt, err := tx.Prepare(sql_add)
    if err != nil {
        panic(err.Error())
    }
    res, err := stmt.Exec(member_id)
    if err != nil {
        panic(err.Error())
    }
    lastId, err := res.LastInsertId()
    if err != nil {
        panic(err.Error())
    }
    rowCnt, err := res.RowsAffected()
    if err != nil {
        panic(err.Error())
    }

    stmt, err = tx.Prepare(sql_mod)
    if err != nil {
        panic(err.Error())
    }
    res, err = stmt.Exec(lastId, variation_id)
    if err != nil {
        panic(err.Error())
    }
    rowCnt, err = res.RowsAffected()
    if err != nil {
        panic(err.Error())
    }

    if rowCnt > 0 {
        stmt, err = tx.Prepare(sql_get)
        if err != nil {
            panic(err.Error())
        }
        rows, err := stmt.Query(lastId)
        if err != nil {
            panic(err.Error())
        }
        seat := GetData(rows);
        PrintData(seat)
        tx.Commit()

        model := BuyViewModel {
            RecentSold: recent,
            SeatId: seat[0]["seat_id"],
            MemberId: member_id,
            Result: true,
        }

        return model
    } else {
        fmt.Print(">> rollback")
        tx.Rollback()

        model := BuyViewModel {
            RecentSold: recent,
            SeatId: "",
            MemberId: "",
            Result: false,
        }
        return model
    }
}

func (p *DbConn) AdminOrder() []map[string]string {
    fmt.Print("DbConn: get /admin/order.csv\n");
    sql := "SELECT order_request.*, stock.seat_id, stock.variation_id, stock.updated_at " +
           "FROM order_request JOIN stock ON order_request.id = stock.order_id " +
           "ORDER BY order_request.id ASC"
    fmt.Printf("  [%s]\n", sql);

    rows, err := p.db_conn.Query(sql)
    if err != nil {
        panic(err.Error())
    }
    defer rows.Close()

    return GetData(rows);
}

func (p *DbConn) Admin() {
    sql_file := "../config/database/initial_data.sql"

    file, err := os.Open(sql_file)
    defer file.Close()
    if err != nil {
        panic(err.Error())
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

        _, err = p.db_conn.Exec(s)
    }
}

////////////////////////////////////////////////////////////////////////////////
func rootHandler(w http.ResponseWriter, r *http.Request) {
    model := db.Index()

    var t = template.Must(template.ParseFiles("template/index.html"))
    if err := t.Execute(w, model); err != nil {
        fmt.Println(err.Error())
    }
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/artist/"):]
    //fmt.Fprintf(w, "artist:\n")
    //fmt.Fprintf(w, "  id: %s\n", id)

    model := db.GetArtists(id)

    var t = template.Must(template.ParseFiles("template/artist.html"))
    if err := t.Execute(w, model); err != nil {
        fmt.Println(err.Error())
    }
}

func ticketHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/ticket/"):]
    //fmt.Fprintf(w, "artist:\n")
    //fmt.Fprintf(w, "  id: %s\n", id)

    model := db.GetTickets(id)

    var t = template.Must(template.ParseFiles("template/ticket.html"))
    if err := t.Execute(w, model); err != nil {
        fmt.Println(err.Error())
    }
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
    variation_id := r.PostFormValue("variation_id")
    member_id := r.PostFormValue("member_id")

    //fmt.Fprintf(w, "buy:\n")
    //fmt.Fprintf(w, "  variation_id: %s\n", variation_id)
    //fmt.Fprintf(w, "  member_id: %s\n", member_id)

    model := db.Buy(variation_id, member_id)

    if model.Result == true {
        var t = template.Must(template.ParseFiles("template/buy_complete.html"))
        if err := t.Execute(w, model); err != nil {
            fmt.Println(err.Error())
        }
    } else {
        var t = template.Must(template.ParseFiles("template/buy_soldout.html"))
        if err := t.Execute(w, model); err != nil {
            fmt.Println(err.Error())
        }
    }
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Print("get /admin\n");

    db.Admin()

    fmt.Fprintf(w, "\n")
}

func adminCsvHandler(w http.ResponseWriter, r *http.Request) {
    data := db.AdminOrder()

    siz := len(data)
    for i := 0; i < siz; i++ {
        for _, v := range data[i] {
            fmt.Fprintf(w, "%s,", v)
        }
        fmt.Fprintf(w, "\n")
    }
}

////////////////////////////////////////////////////////////////////////////////
func main() {
    db.Init()

    http.HandleFunc("/artist/", artistHandler);
    http.HandleFunc("/ticket/", ticketHandler);
    http.HandleFunc("/buy/", buyHandler);
    http.HandleFunc("/admin/", adminHandler);
    http.HandleFunc("/admin/order.csv", adminCsvHandler);
    http.HandleFunc("/", rootHandler);

    http.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })
    http.HandleFunc("/css/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })
    http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })

    http.ListenAndServe(":5000", nil)

    //defer db_conn.Close()
}

