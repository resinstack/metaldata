package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-hclog"
)

// Server is a wrapper around an Echo and a logger.
type Server struct {
	l hclog.Logger

	r chi.Router
	n *http.Server

	source InfoSource
}

// InfoSource is a provider of machine metadata.
type InfoSource interface {
	GetMachineInfo(string, string) (string, error)
	GetUserData(string) (string, error)
}
