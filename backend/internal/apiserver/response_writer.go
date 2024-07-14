package apiserver

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.code = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
