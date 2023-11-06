package utils

import (
	"fmt"
	"net/http"
)

type ResponseWriter struct {
	Status int
	Body   string
	Time   int64
	http.ResponseWriter
}

func NewResponseWriter(res http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: res}
}

func (res *ResponseWriter) WriteHeader(code int) {
	res.Status = code
	res.ResponseWriter.WriteHeader(code)
}

func (res *ResponseWriter) Write(body []byte) (int, error) {
	res.Body = string(body)
	return res.ResponseWriter.Write(body)
}

func (res *ResponseWriter) String() string {
	out := fmt.Sprintf("status %d (took %dms)", res.Status, res.Time)
	if res.Body != "" {
		out = fmt.Sprintf("%s\n\tresponse: %s", out, res.Body)
	}
	return out
}
