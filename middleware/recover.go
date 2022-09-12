package middleware

import (
	"log"
	"net/http"
)

type Recover struct {
	next http.Handler
}

func NewRecover(next http.Handler) *Recover {
	return &Recover{next: next}
}

func (r *Recover) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered: %v", r)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	r.next.ServeHTTP(w, req)
}
