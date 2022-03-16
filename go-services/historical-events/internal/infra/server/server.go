package server

type grpc interface {
	Start()
}

type http interface {
	Start()
}

// Server describes the server structure alongside it's dependencies.
type Server struct {
	grpc grpc
	http http
}

// New creates the new server
func New(grpc grpc, http http) *Server {
	return &Server{
		grpc: grpc,
		http: http,
	}
}

// StartGRPC starts grpc server
func (s *Server) StartGRPC() {
	s.grpc.Start()
}

// StartHTTP  starts http server
func (s *Server) StartHTTP() {
	s.http.Start()
}
