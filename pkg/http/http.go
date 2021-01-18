package http

import (
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
)

// New returns an HTTP server
func New(parent hclog.Logger) *Server {
	x := &Server{
		l:    parent.Named("http"),
		Echo: echo.New(),
	}

	x.GET("/", x.peerInfo)

	// Handle most metadata keys
	x.GET("/get/meta/:key", x.getMetaData)

	return x
}

func (s *Server) peerInfo(c echo.Context) error {
	res, err := s.getPeerID(c.Request())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, res)
}

func (s *Server) getMetaData(c echo.Context) error {
	return s.handleMetadataRequest(c, c.Param("key"))
}

func (s *Server) handleMetadataRequest(c echo.Context, key string) error {
	hwaddr, err := s.getPeerID(c.Request())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	value, err := s.source.GetMachineInfo(hwaddr, key)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, value)
}
