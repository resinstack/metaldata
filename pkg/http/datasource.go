package http

import (
	"net"
	"net/http"
	"strings"

	"github.com/mostlygeek/arp"
)

// SetSource set up the source that will be served by this instance.
func (s *Server) SetSource(i InfoSource) {
	s.source = i
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
