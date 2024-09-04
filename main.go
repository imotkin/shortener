package main

import "log"

func main() {
	srv := NewServer("db.sqlite3")
	if err := srv.Start("0.0.0.0:5000"); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
