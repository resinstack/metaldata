package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-hclog"
)

// New returns an HTTP server
func New(opts ...Option) *Server {
	x := &Server{
		l: hclog.NewNullLogger(),
		r: chi.NewRouter(),
		n: &http.Server{},
	}

	for _, o := range opts {
		o(x)
	}

	x.r.Get("/", x.peerInfo)
	x.r.Get("/get/meta/{key}", x.getMetaData)
	x.r.Get("/get/user", x.getUserData)
	return x
}

// Serve binds and serves http on the bound socket.  An error will be
// returned if the server cannot initialize.
func (s *Server) Serve(bind string) error {
	s.l.Info("HTTP is starting")
	s.n.Addr = bind
	s.n.Handler = s.r
	return s.n.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.n.Shutdown(ctx)
}

func (s *Server) jsonError(w http.ResponseWriter, r *http.Request, rc int, err error) {
	w.WriteHeader(rc)
	enc := json.NewEncoder(w)
	err = enc.Encode(struct {
		Error error
	}{
		Error: err,
	})
	if err != nil {
		fmt.Fprintf(w, "Error writing json error response: %v", err)
	}
}

func (s *Server) peerInfo(w http.ResponseWriter, r *http.Request) {
	res, err := s.getPeerID(r)
	if err != nil {
		s.jsonError(w, r, http.StatusInternalServerError, err)
		return
	}
	fmt.Fprint(w, res)
}

func (s *Server) getMetaData(w http.ResponseWriter, r *http.Request) {
	s.handleMetadataRequest(w, r, chi.URLParam(r, "key"))
}

func (s *Server) getUserData(w http.ResponseWriter, r *http.Request) {
	hwaddr, err := s.getPeerID(r)
	if err != nil {
		s.l.Warn("Error getting peer ID", "error", err)
		s.jsonError(w, r, http.StatusInternalServerError, err)
		return
	}

	value, err := s.source.GetUserData(hwaddr)
	if err != nil {
		s.l.Warn("Error loading user data", "error", err)
		s.jsonError(w, r, http.StatusInternalServerError, err)
		return
	}
	fmt.Fprint(w, value)
}

func (s *Server) handleMetadataRequest(w http.ResponseWriter, r *http.Request, key string) {
	hwaddr, err := s.getPeerID(r)
	if err != nil {
		s.l.Warn("Error getting peer ID", "error", err)
		s.jsonError(w, r, http.StatusInternalServerError, err)
		return
	}

	value, err := s.source.GetMachineInfo(hwaddr, key)
	if err != nil {
		s.l.Warn("Error loading metadata", "error", err, "key", key)
		s.jsonError(w, r, http.StatusInternalServerError, err)
		return
	}
	fmt.Fprint(w, value)
}
