package turn

import (
	"github.com/pixelbender/go-stun/mux"
	"github.com/pixelbender/go-stun/stun"
	"net"
	"time"
)

type Server struct {
	mux.Mux
	LifeTime time.Duration
}

func NewServer() *Server {
	srv := &Server{
		LifeTime: time.Second,
	}
	srv.Handle(srv.ServeMux)
	return srv
}

func (srv *Server) ServeSTUN(rw stun.ResponseWriter, r *stun.Message) {
	switch r.Method {
	case MethodAllocate:
		rw.WriteMessage(&stun.Message{
			Method: r.Method | stun.TypeResponse,
			Attributes: stun.Attributes{
				stun.AttrXorMappedAddress: rw.RemoteAddr(),
				AttrXorRelayedAddress:     &stun.Addr{IP: net.ParseIP("127.0.0.1"), Port: 1000},
				AttrLifeTime:              uint32(srv.LifeTime / time.Second),
			},
		})
	default:
		srv.Handler.ServeSTUN(rw, r)
	}
}

type allocation struct {
	from *Conn
	to   net.Conn
}
