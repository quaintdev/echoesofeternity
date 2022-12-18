package main

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "quotes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	http.HandleFunc("/", homepage)
	http.HandleFunc("/qotd", quoteOfTheDay)
	http.HandleFunc("/quotes", listQuotes)
	http.ListenAndServe(":8080", nil)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func quoteOfTheDay(w http.ResponseWriter, r *http.Request) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM quotes").Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var idx = rand.Intn(count)
	var quote Quote
	err = db.QueryRow("SELECT author, quote FROM quotes WHERE id = ?", idx).Scan(&quote.Author, &quote.Quote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(quote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func listQuotes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, author, quote FROM quotes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var quotes []Quote
	for rows.Next() {
		var quote Quote
		err := rows.Scan(&quote.ID, &quote.Author, &quote.Quote)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		quotes = append(quotes, quote)
	}

	json, err := json.Marshal(quotes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
