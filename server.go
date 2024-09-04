package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"text/template"
	"time"

	_ "modernc.org/sqlite"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lenID   = 7
)

var (
	mainTemplate    = template.Must(template.ParseFiles("./static/index.html"))
	errorTemplate   = template.Must(template.ParseFiles("./static/error.html"))
	statsTemplate   = template.Must(template.ParseFiles("./static/stats.html"))
	shortenTemplate = template.Must(template.ParseFiles("./static/shorten.html"))
)

func randomID() string {
	b := make([]byte, lenID)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type Server struct {
	mux *http.ServeMux
	db  *Database
}

func NewServer(pathDB string) *Server {
	return &Server{mux: http.NewServeMux(), db: NewDatabase(pathDB)}
}

func (s *Server) mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		s.Template(mainTemplate, w, nil)
	} else {
		s.ErrorCode(w, http.StatusNotFound)
	}
}

func (s *Server) linkHandler(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("URL")

	URL, err := s.db.Get(ID)
	if err != nil {
		s.ErrorCode(w, http.StatusNotFound)
		return
	}

	IP := cmp.Or(r.Header.Get("X-Forwarded-For"), r.RemoteAddr)

	loc, err := FindLocation(IP)
	if err != nil {
		log.Printf("find location: %v", err)
	}

	err = s.db.UpdateStats(URL.Original, loc, IP)
	if err != nil {
		log.Printf("add link view: %v", err)
	}

	var redirect string

	if strings.HasPrefix(URL.Original, "http") {
		redirect = URL.Original
	} else {
		redirect = "http://" + URL.Original
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}

func (s *Server) shortenHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.ErrorCode(w, http.StatusBadRequest)
		log.Printf("parse URL value: %v", err)
		return
	}

	URL := r.PostForm.Get("link")

	link, err := s.db.Get(URL)
	if err == nil { // if URL was found
		s.Template(shortenTemplate, w, map[string]string{
			"MainURL":  fmt.Sprintf("http://%s/%s", r.Host, link.Shortened),
			"StatsURL": fmt.Sprintf("http://%s/stats/%s", r.Host, link.Shortened),
		})
		return
	}

	ID := randomID()

	if err := s.db.Add(URL, ID); err != nil {
		s.ErrorCode(w, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	s.Template(shortenTemplate, w, map[string]string{
		"MainURL":  fmt.Sprintf("http://%s/%s", r.Host, ID),
		"StatsURL": fmt.Sprintf("http://%s/stats/%s", r.Host, ID),
	})
}

func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("URL")

	URL, err := s.db.Get(ID)
	if err != nil {
		s.ErrorCode(w, http.StatusNotFound)
		return
	}

	data := map[string]any{
		"ID":         ID,
		"MainURL":    fmt.Sprintf("http://%s/%s", r.Host, ID),
		"Views":      URL.Views,
		"LatestView": URL.LatestView,
		"CreatedAt":  URL.CreatedAt.Format("02-01-2006 15:04:05"),
	}

	s.Template(statsTemplate, w, data)
}

func (s *Server) viewsHandler(w http.ResponseWriter, r *http.Request) {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			ID := r.PathValue("URL")

			URL, err := s.db.Get(ID)
			if err != nil {
				s.ErrorCode(w, http.StatusNotFound)
				return
			}

			data, err := json.Marshal(fmt.Sprintf(`{"views": %d}`, URL.Views))
			if err != nil {
				s.ErrorCode(w, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			fmt.Fprintf(w, "event: %s\n", "views-update")
			fmt.Fprintf(w, "data: %s\n\n", string(data))
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) ErrorCode(w http.ResponseWriter, code int) {
	s.Template(errorTemplate, w, fmt.Sprintf("%d %s", code, http.StatusText(code)))
}

func (s *Server) ErrorMessage(w http.ResponseWriter, msg string) {
	s.Template(errorTemplate, w, msg)
}

func (s *Server) Template(tmpl *template.Template, w http.ResponseWriter, data any) {
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) Start(addr string) error {
	fs := http.FileServer(http.Dir("./static"))
	s.mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	s.mux.HandleFunc("GET /stats/{URL}", s.statsHandler)
	s.mux.HandleFunc("GET /views/{URL}", s.viewsHandler)
	s.mux.HandleFunc("GET /{URL}", s.linkHandler)
	s.mux.HandleFunc("GET /", s.mainHandler)

	s.mux.HandleFunc("POST /shorten", s.shortenHandler)

	return http.ListenAndServe(addr, s.mux)
}
