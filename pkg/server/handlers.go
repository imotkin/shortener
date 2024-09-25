package server

import (
	"encoding/json"
	"fmt"
	"github.com/imotkin/shortener/pkg/api"
	"log"
	"net/http"
	"strings"
	"time"
)

func (s *Server) mainHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			s.Template(mainTemplate, w, nil)
		} else {
			s.ErrorCode(w, http.StatusNotFound)
		}
	})
}

func (s *Server) linkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID := r.PathValue("URL")

		URL, err := s.db.Get(ID)
		if err != nil {
			s.ErrorCode(w, http.StatusNotFound)
			return
		}

		var (
			loc api.Response
			IP  = ParseIP(r)
		)

		if IP != "0.0.0.0" {
			loc, err = api.FindLocation(IP)
			if err != nil {
				log.Printf("Find location: %v", err)
			}
		}

		err = s.db.UpdateStats(URL.Original, loc, IP)
		if err != nil {
			log.Printf("Add link view: %v", err)
		}

		var redirect string

		if strings.HasPrefix(URL.Original, "http") {
			redirect = URL.Original
		} else {
			redirect = "http://" + URL.Original
		}

		http.Redirect(w, r, redirect, http.StatusFound)
	})
}

func (s *Server) shortenHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			s.ErrorCode(w, http.StatusBadRequest)
			log.Printf("Parse URL value: %v", err)
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

		ID := RandomID()

		if err := s.db.Add(URL, ID); err != nil {
			s.ErrorCode(w, http.StatusInternalServerError)
			log.Printf("Add URL: %v", err)
			return
		}

		s.Template(shortenTemplate, w, map[string]string{
			"MainURL":  fmt.Sprintf("http://%s/%s", r.Host, ID),
			"StatsURL": fmt.Sprintf("http://%s/stats/%s", r.Host, ID),
		})
	})
}

func (s *Server) statsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func (s *Server) viewsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}
