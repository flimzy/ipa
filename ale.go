package ale

import (
	"context"
	"time"

	"github.com/flimzy/log"
	"github.com/julienschmidt/httprouter"
	"github.com/tylerb/graceful"
)

// Logger is an interface to a minimal logger, such as the default *log.Logger, or my
// preferred github.com/flimzy/log.Logger.
type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

// RouteOptions provides route-specific configuration
type RouteOptions struct {
	Template string
}

// Debugger is an interface to a debugger, such as github.com/flimzy/log.Logger
type Debugger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Debugln(...interface{})
}

// Timeout defines the default time to wait before killing active connections on shutdown or restart.
const Timeout = 10 * time.Second

// Server represents an Ale server instance.
type Server struct {
	// Timeout is the duration to wait before killing active requests when stopping the server
	Timeout time.Duration
	// router is an instance of julienschmidt/httprouter
	router *httprouter.Router
	// Context is the master context for this server instance
	Context     Context
	httpServer  *graceful.Server
	httpsServer *graceful.Server
	envPrefix   string
	err         error
}

// New returns a new Ale server instance.
func New() *Server {
	s := &Server{
		Timeout: Timeout,
		Context: context.Background(),
		router:  httprouter.New(),
	}
	return s
}

// Run initializes the web server instance
func (s *Server) Run() error {
	httpAddr := s.GetConf(ConfHTTPBind)
	httpsAddr := s.GetConf(ConfHTTPSBind)

	log.Debugf("Run(). httpAddr = %s, httpsAddr = %s", httpAddr, httpsAddr)

	if httpAddr != "" && httpsAddr != "" {
		return s.serveBoth(httpsAddr, httpAddr)
	}
	if httpAddr != "" {
		return s.serveHTTP(httpAddr)
	}
	return s.serveHTTPS(httpsAddr)
}
