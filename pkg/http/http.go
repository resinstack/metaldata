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
	x.GET("/latest/meta-data/:key", x.getMetaData)

	// Handle keys that are nested or otherwise need special
	// handling.
	x.GET("/latest/meta-data/placement/availability-zone", x.getAvailabilityZone)
	x.GET("/latest/meta-data/public-keys/0/openssh-key", x.getSSHKey)
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

func (s *Server) getAvailabilityZone(c echo.Context) error {
	return s.handleMetadataRequest(c, "availability-zone")
}

func (s *Server) getSSHKey(c echo.Context) error {
	return s.handleMetadataRequest(c, "ssh-key")
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
