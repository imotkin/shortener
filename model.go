package main

import (
	"time"
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
	IP, Country, City, VisitTime string
	// VisitTime         time.Time
}
