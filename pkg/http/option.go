package http

import (
	"github.com/hashicorp/go-hclog"
)

// WithLogger configures the logger for the webserver.
func WithLogger(l hclog.Logger) Option {
	return func(s *Server) {
		s.l = l.Named("http")
	}
}

// WithInfoSource specifies a lookup service for information.
func WithInfoSource(is InfoSource) Option {
	return func(s *Server) {
		s.source = is
	}
}
