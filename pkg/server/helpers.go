package server

import (
	"cmp"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lenID   = 7
)

func NewTemplate(path string) *template.Template {
	return template.Must(template.ParseFiles(path))
}

func RandomID() string {
	b := make([]byte, lenID)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func ParseIP(r *http.Request) string {
	host := r.RemoteAddr

	if strings.HasPrefix(host, "127.0.0.1") || strings.HasPrefix(host, "::1") {
		return "0.0.0.0"
	}

	return cmp.Or(r.Header.Get("X-Forwarded-For"), r.RemoteAddr)
}
