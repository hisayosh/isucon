package main

import (
    "fmt"
    "net/http"
    "html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Artist struct {
	Id string
	Name string
}

type Recent struct {
	Seat_id string
	V_name string
	T_name string
	A_name string
}

type Info struct {
	At []Artist
	Rt []Recent
}

func top_handler(w http.ResponseWriter, r *http.Request) {


	db, err := sql.Open("mysql", "isucon2app:isunageruna@/isucon2")
	checkErr(err)
	artists, err := db.Query("SELECT * FROM artist")

	artists_slice := make([]Artist, 0, 10)

	for artists.Next() {
		var a Artist
		err = artists.Scan(&a.Id, &a.Name)
		checkErr(err)
		artists_slice = append(artists_slice, a) 
	}

	checkErr(err)
	var recent_slice = get_recent_sold(db)

	info := Info{
		At: artists_slice,
		Rt: recent_slice,
	}

	var t = template.Must(template.ParseFiles("layout.html", "index.html"))
	if err := t.ExecuteTemplate(w, "layout", info); err != nil {
//	if err := t.Execute(w, artists_slice); err != nil {
		fmt.Println(err.Error())
	}

}

func get_recent_sold(db *sql.DB) []Recent {
	
	recent_sold, err := db.Query("SELECT stock.seat_id, variation.name AS v_name, ticket.name AS t_name, artist.name AS a_name FROM stock JOIN variation ON stock.variation_id = variation.id JOIN ticket ON variation.ticket_id = ticket.id JOIN artist ON ticket.artist_id = artist.id WHERE order_id IS NOT NULL ORDER BY order_id DESC LIMIT 10")
	checkErr(err)

	recent_slice := make([]Recent, 0, 100)

	for recent_sold.Next() {

		var r Recent
		err = recent_sold.Scan(&r.Seat_id, &r.V_name, &r.T_name, &r.A_name)
		checkErr(err)

		recent_slice = append(recent_slice, r)
	}

	return recent_slice
}

func main() {
    http.HandleFunc("/", top_handler)
	http.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/css/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
    panic(http.ListenAndServe(":5555", nil))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
