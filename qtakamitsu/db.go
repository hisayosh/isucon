package main

import (
  "os"
  "bufio"
  "io"
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

type DbConn struct {
    db_conn *sql.DB
}

func (p *DbConn) Init() {
    fmt.Print("DbConn: init\n");

    db_string := fmt.Sprintf("%s:%s@/%s", conf.Username, conf.Password, conf.Dbname)
    db_conn, err := sql.Open("mysql", db_string)
    if err != nil {
        fmt.Print("error: open failed.\n");
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
    Artist []map[string]string
}

func (p *DbConn) Index() IndexViewModel {
    fmt.Print("DbConn: get /\n");
    sql := "SELECT * FROM artist ORDER BY id"
    fmt.Printf("  [%s]\n", sql);

    rows, err := p.db_conn.Query(sql)
    if err != nil {
        panic(err.Error())
    }
    defer rows.Close()

    recent := db.RecentSold()
    artist := GetData(rows)

    model := IndexViewModel {
        RecentSold: recent,
        Artist: artist,
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
    Artist []map[string]string
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

    model := ArtistViewModel {
        Artist: artists,
        Tickets: tickets,
    }

    return model
}

type TicketViewModel struct {
    Tickets []map[string]string
    Variations []map[string]string
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
        fmt.Println("Error: query failed");
        panic(err.Error())
    }
    rows, err := stmt.Query(ticket_id)
    if err != nil {
        fmt.Println("Error: query failed");
        panic(err.Error())
    }
    tickets := GetData(rows);
    //PrintData(ticket)

    stmt, err = p.db_conn.Prepare(sql_variations)
    if err != nil {
        fmt.Println("Error: query failed");
        panic(err.Error())
    }
    rows, err = stmt.Query(tickets[0]["id"])
    if err != nil {
        fmt.Println("Error: query failed");
        panic(err.Error())
    }
    defer rows.Close()
    variations := GetData(rows);
    //PrintData(variations)

    siz := len(variations)
    for i := 0; i < siz; i++ {
        stmt, err = p.db_conn.Prepare(sql_variation_stock)
        if err != nil {
            fmt.Println("Error: query failed");
            panic(err.Error())
        }
        rows, err = stmt.Query(variations[i]["id"])
        if err != nil {
            fmt.Println("Error: query failed");
            panic(err.Error())
        }
        //stack := GetData(rows);
        //PrintData(stack)

        stmt, err = p.db_conn.Prepare(sql_variation_vacancy)
        if err != nil {
            fmt.Println("Error: query failed");
            panic(err.Error())
        }
        rows, err = stmt.Query(variations[i]["id"])
        if err != nil {
            fmt.Println("Error: query failed");
            panic(err.Error())
        }
        //num := GetData(rows);
        //PrintData(num)
    }

    model := TicketViewModel {
        Tickets: tickets,
        Variations: variations,
    }

    return model
}

type BuyViewModel struct {
    SeatId string
    MemberId string
}

func (p *DbConn) Buy(variation_id string, member_id string) BuyViewModel {
    fmt.Printf("DbConn: post /buy (%s:%s)\n", variation_id, member_id)

    sql_add := "INSERT INTO order_request (member_id) VALUES (?)"
    sql_mod := "UPDATE stock SET order_id = ? WHERE variation_id = ? AND order_id IS NULL ORDER BY RAND() LIMIT 0"
    sql_get := "SELECT seat_id FROM stock WHERE order_id = ? LIMIT 1', $order_id"

    fmt.Printf("  [%s]\n", sql_add);
    fmt.Printf("  [%s]\n", sql_mod);
    fmt.Printf("  [%s]\n", sql_get);

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
            SeatId: "",
            MemberId: "",
        }

        return model
    } else {
        fmt.Print(">> rollback")
        tx.Rollback()

        model := BuyViewModel { "", "" }
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



var db DbConn

func main() {

    db.Init()

/*
    data := db.RecentSold()
    for i := 0; i < len(data); i++ {
        for k, v := range data[i] {
            fmt.Printf("%s=%s, ", k, v)
        }
        fmt.Print("\n")
    }
*/

/*
    data = db.Index()
    for i := 0; i < len(data); i++ {
        for k, v := range data[i] {
            fmt.Printf("%s=%s, ", k, v)
        }
        fmt.Print("\n")
    }
*/

/*
    db.GetArtists("1")
*/

    db.GetTickets("2")
/*
*/

/*
    db.Buy("3", "100")
*/
/*
    data := db.AdminOrder()
    PrintData(data)
*/

    db.Admin()

    db.Close()
}

