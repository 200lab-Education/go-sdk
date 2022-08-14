package httpserver

import (
	"net"
	"net/http"
	"time"
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	_ = tc.SetKeepAlive(true)
	_ = tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

type myHttpServer struct {
	http.Server
}

func (srv *myHttpServer) Serve(lis net.Listener) error {
	return srv.Server.Serve(tcpKeepAliveListener{lis.(*net.TCPListener)})
}

func (srv *myHttpServer) ServeTLS(lis net.Listener, certFile, keyFile string) error {
	return srv.Server.ServeTLS(tcpKeepAliveListener{lis.(*net.TCPListener)}, certFile, keyFile)
}
