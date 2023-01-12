package model

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, //1MB,
	}
	return s.httpServer.ListenAndServe()
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)

}

type TimeFormat struct {
	TimeStamp int64
	Text      string
}

func TimeFormatter(t int64) TimeFormat {
	var result TimeFormat
	result.TimeStamp = t
	if t == 0 {
		result.Text = "empty"
	} else {
		unix := time.Unix(t, 0)
		result.Text = unix.Format(time.UnixDate)
	}
	return result

}
