package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"

	"github.com/imotkin/shortener/internal/ip"
	"github.com/imotkin/shortener/internal/migrations"
)

type URL struct {
	ID         int
	Original   string
	Shortened  string
	CreatedAt  time.Time
	Views      int
	LatestView string
}

type Stats struct {
	IP        string
	Country   string
	City      string
	VisitTime string
}

type Database struct {
	conn *sql.DB
}

func New(path string) *Database {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("Connect to database: %v", err)
	}

	log.Printf("Load database: %s", path)

	err = migrations.RunMigrations(db)
	if err != nil {
		log.Fatalf("Run migrations: %v", err)
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

func (db *Database) UpdateStats(link string, loc ip.Response, IP string) error {
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
		  FROM stats WHERE link_id = (SELECT id FROM links WHERE shortened = ?)`, ID)
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
