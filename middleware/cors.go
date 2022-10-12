package middleware

import (
	"log"
	"net/http"
)

type Cors struct {
	next http.Handler
}

func NewCors(next http.Handler) *Cors {
	return &Cors{next: next}
}

func (c *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("r method %v", r.Method)
	w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

	c.next.ServeHTTP(w, r)
}
