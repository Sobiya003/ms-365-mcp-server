package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/softeria/ms-365-mcp-server/internal/auth"
	"github.com/softeria/ms-365-mcp-server/internal/cli"
)

type Server struct {
	auth    *auth.Manager
	opts    cli.Options
	version string
}

func New(a *auth.Manager, o cli.Options, version string) *Server {
	return &Server{auth: a, opts: o, version: version}
}

func (s *Server) Start() error {
	if s.opts.HTTP == "" {
		return s.runStdio()
	}
	return s.runHTTP(s.opts.HTTP)
}

func (s *Server) runStdio() error {
	resp := map[string]any{
		"name":      "ms-365-mcp-server",
		"version":   s.version,
		"cloud":     s.opts.Cloud,
		"readOnly":  s.opts.ReadOnly,
		"orgMode":   s.opts.OrgMode,
		"discovery": s.opts.Discovery,
	}
	return json.NewEncoder(os.Stdout).Encode(resp)
}

func (s *Server) runHTTP(bind string) error {
	host, port := parseHostPort(bind)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	mux.HandleFunc("/.well-known/oauth-authorization-server", func(w http.ResponseWriter, r *http.Request) {
		origin := fmt.Sprintf("http://%s", r.Host)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"issuer":                 origin,
			"authorization_endpoint": origin + "/authorize",
			"token_endpoint":         origin + "/token",
		})
	})
	mux.HandleFunc("/mcp", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "Go MCP endpoint placeholder"})
	})

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("HTTP mode listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func parseHostPort(input string) (string, string) {
	if input == "" {
		return "0.0.0.0", "3000"
	}
	if strings.Contains(input, ":") {
		parts := strings.SplitN(input, ":", 2)
		h := parts[0]
		p := parts[1]
		if h == "" {
			h = "0.0.0.0"
		}
		if p == "" {
			p = "3000"
		}
		return h, p
	}
	return "0.0.0.0", input
}
