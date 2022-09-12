package middleware

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

type Logger struct {
	next http.Handler
}

func NewLogger(next http.Handler) *Logger {
	return &Logger{next: next}
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := &loggerResponseWriter{ResponseWriter: w}

	start := time.Now()
	defer func() {
		log.Printf("[%d] %s %s %s\nresp: %s\n", response.code, r.Method, r.URL.Path, time.Since(start), response.resp.String())
	}()

	l.next.ServeHTTP(response, r)
}

type loggerResponseWriter struct {
	http.ResponseWriter
	code int
	resp bytes.Buffer
}

func (l *loggerResponseWriter) WriteHeader(code int) {
	if l.code != 0 {
		l.code = code
		l.ResponseWriter.WriteHeader(code)
	}
}

func (l *loggerResponseWriter) Write(p []byte) (int, error) {
	l.WriteHeader(http.StatusOK)
	l.resp.Write(p)
	return l.ResponseWriter.Write(p)
}
