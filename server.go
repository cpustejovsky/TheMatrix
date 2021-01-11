package matrix

import (
	"errors"
	"io"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	// The address the server will listen on for incoming tcp connections
	ListenAddress string
}

type server struct {
	conf     ServerConfig
	listener net.Listener
	wg       sync.WaitGroup
}

func SpawnServer(conf ServerConfig) (io.Closer, error) {
	s := server{conf: conf}
	var err error

	// Attempt to listen at the address provided
	s.listener, err = net.Listen("tcp", s.conf.ListenAddress)
	if err != nil {
		return nil, err
	}

	// Main accept loop
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		var id uint64

		for {
			// Wait for new connections
			conn, err := s.listener.Accept()
			if err != nil {
				// Listener temporary errors
				if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
					logrus.WithError(err).Warn("temporary listener error")
					continue
				}

				// Listener was closed
				if errors.Unwrap(err).Error() == "use of closed network connection" {
					return
				}

				logrus.WithError(err).Error("network error; closing")
				return
			}

			// Add to the work group and increment our id count
			s.wg.Add(1)
			id++

			// Spawn a new go routine to handle the connection
			go func(id uint64) {
				defer s.wg.Done()
				if err := s.handleConn(id, conn); err != nil {
					logrus.WithError(err).Error("internal error; internal error")
				}
			}(id)
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

func (s *server) handleConn(id uint64, conn net.Conn) error {
	logrus.WithField("id", id).Infof("New connection (%s)", conn.RemoteAddr())
	if _, err := conn.Write([]byte("Welcome to the matrix\r\n")); err != nil {
		return err
	}
	return conn.Close()
}
