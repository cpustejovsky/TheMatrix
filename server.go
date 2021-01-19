package matrix

import (
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/mailgun/holster/v3/setter"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	// The address the server will listen on for incoming tcp connections
	ListenAddress string

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. Does not let Handlers make
	// per-request read timeouts
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	ReadHeaderTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeout time.Duration

	Log logrus.FieldLogger
}

type server struct {
	conf     ServerConfig
	listener net.Listener
	wg       sync.WaitGroup
	srv      http.Server
}

func SpawnServer(conf ServerConfig) (io.Closer, error) {
	setter.SetDefault(&conf.Log, logrus.WithField("category", "server"))
	s := server{conf: conf}
	var err error

	// Attempt to listen at the address provided
	s.listener, err = net.Listen("tcp", s.conf.ListenAddress)
	if err != nil {
		return nil, err
	}

	s.srv = http.Server{
		Handler:      &s,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		IdleTimeout:  conf.IdleTimeout,
		// TODO: Ensure all logging is done under logrus
		//ErrorLog:        nil,
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.srv.Serve(s.listener); err != nil {
			if err != http.ErrServerClosed {
				s.conf.Log.WithError(err).Error("while running http.Serve()")
			}
		}
	}()
	return &s, nil
}

func (s *server) Close() error {
	defer s.wg.Wait()
	if err := s.listener.Close(); err != nil {
		return err
	}
	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.conf.Log.Infof("New connection from '%s'", r.RemoteAddr)
	if _, err := w.Write([]byte("Welcome to the matrix\r\n")); err != nil {
		s.conf.Log.WithError(err).Error("while writing to client")
	}
}
