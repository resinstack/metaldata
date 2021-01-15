package http

import (
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"github.com/mostlygeek/arp"
)

// New returns an HTTP server
func New(parent hclog.Logger) *Server {
	x := &Server{
		l:    parent.Named("http"),
		Echo: echo.New(),
	}

	x.GET("/", x.peerInfo)

	return x
}

func (s *Server) peerInfo(c echo.Context) error {
	res, err := s.getPeerID(c.Request())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, res)
}

func (s *Server) getPeerID(r *http.Request) (string, error) {
	addr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if addr == "127.0.0.1" || addr == "::1" {
		// localhost is handled specially to allow testing to
		// be done easier.  It is always returned with a
		// hwaddr of all zeros.
		return "00:00:00:00:00:00", nil
	}

	// Slightly wasteful, but in the general case this is probably
	// the only time we'll hear from this machine as it slurps a
	// bunch of data, so its worth just refreshing the in-memory
	// cache.
	arp.CacheUpdate()
	hwaddr := strings.ToUpper(arp.Search(addr))
	s.l.Trace("getPeerID", "address", addr, "hwaddr", hwaddr)
	return hwaddr, nil
}
