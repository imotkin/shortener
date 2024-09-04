package main

import (
	"database/sql"
	"fmt"
	"log"
)

type LinkRepository interface {
	Add(URL, ID string) error
	Get(ID string, counter bool) (URL, error)
}

type Database struct {
	conn *sql.DB
}

func NewDatabase(path string) *Database {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS links (
	            		id INTEGER PRIMARY KEY AUTOINCREMENT, 
						original TEXT NOT NULL,
						shortened TEXT NOT NULL,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						views INTEGER NOT NULL DEFAULT(0));
						
					  CREATE TABLE IF NOT EXISTS stats (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						link_id INTEGER NOT NULL,
						ip TEXT,
						visit_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						country TEXT,
						region TEXT,
						city TEXT,
						FOREIGN KEY (link_id) REFERENCES links(id))`)
	if err != nil {
		log.Fatalf("create database tables: %v", err)
	}

	return &Database{conn: db}
}

func (db *Database) Add(URL, ID string) error {
	_, err := db.conn.Exec(`INSERT INTO links (original, shortened) VALUES (?, ?)`, URL, ID)
	if err != nil {
		return fmt.Errorf("add new URL to database: %w", err)
	}

	return nil
}

func (db *Database) Get(ID string) (u URL, err error) {
	row := db.conn.QueryRow(`
		SELECT original, shortened, created_at, 
	           (SELECT COUNT (*) FROM stats WHERE link_id = l.id) AS views,
	           (SELECT COALESCE((SELECT country || ', ' || city || ', ' || strftime('%m-%d-%Y %H:%M:%S', visit_time) || ' GMT'
                  FROM stats 
                 WHERE link_id = l.id 
                 ORDER BY visit_time DESC LIMIT 1), "")) AS last_visit
          FROM links l
         WHERE shortened = ? OR original = ?`, ID, ID)

	err = row.Scan(&u.Original, &u.Shortened, &u.CreatedAt, &u.Views, &u.LatestView)
	if err != nil {
		return URL{}, fmt.Errorf("get URL in database: %w", err)
	}

	return
}

func (db *Database) UpdateStats(link string, loc Response, IP string) error {
	_, err := db.conn.Exec(`
		INSERT INTO stats (link_id, ip, country, region, city) 
     	VALUES ((SELECT id FROM links WHERE original = ?), ?, ?, ?, ?)`,
		link, IP, loc.Country, loc.Region, loc.City,
	)
	if err != nil {
		return fmt.Errorf("add new stats to database: %w", err)
	}

	return nil
}

func (db *Database) Stats(ID string) (stats []Stats, err error) {
	rows, err := db.conn.Query(`
		SELECT ip, country, city, strftime('%m-%d-%Y %H:%M:%S', visit_time)
		  FROM stats WHERE link_id = (
		  	SELECT id 
			  FROM links 
			 WHERE shortened = ?)
		`, ID,
	)
	if err != nil {
		return nil, fmt.Errorf("get stats from database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s Stats
		err = rows.Scan(&s.IP, &s.Country, &s.City, &s.VisitTime)
		stats = append(stats, s)
	}

	if err != nil {
		return nil, fmt.Errorf("get stats from database: %w", err)
	}

	return
}
