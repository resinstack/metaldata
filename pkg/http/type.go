package http

import (
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
)

// Server is a wrapper around an Echo and a logger.
type Server struct {
	*echo.Echo
	l hclog.Logger

	source InfoSource
}

// InfoSource is a provider of machine metadata.
type InfoSource interface {
	GetMachineInfo(string, string) (string, error)
}
