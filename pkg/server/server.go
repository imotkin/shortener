package server

import (
	"fmt"
	"github.com/imotkin/shortener/pkg/database"
	"github.com/imotkin/shortener/pkg/middleware"
	"html/template"
	"log"
	"net/http"
)

var (
	mainTemplate    = NewTemplate("static/index.html")
	errorTemplate   = NewTemplate("static/error.html")
	statsTemplate   = NewTemplate("static/stats.html")
	shortenTemplate = NewTemplate("static/shorten.html")
)

type Server struct {
	mux *http.ServeMux
	db  *database.Database
}

func New(pathDB string) *Server {
	return &Server{mux: http.NewServeMux(), db: database.New(pathDB)}
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

	s.mux.Handle("GET /stats/{URL}", middleware.Logger(s.statsHandler()))
	s.mux.Handle("GET /views/{URL}", middleware.Logger(s.viewsHandler()))
	s.mux.Handle("GET /{URL}", middleware.Logger(s.linkHandler()))
	s.mux.Handle("GET /", middleware.Logger(s.mainHandler()))

	s.mux.Handle("POST /shorten", middleware.Logger(s.shortenHandler()))

	log.Printf("Listening at: http://%s", addr)
	return http.ListenAndServe(addr, s.mux)
}
