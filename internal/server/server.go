package server

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/GolubAlexander/grpc-logs/internal/services/generator"
)

type Server struct {
	srv *grpc.Server
	svc Generator
	UnimplementedLoggerServer
}

type Generator interface {
	Logs() chan generator.Message
}

func New(svc Generator) *Server {
	grpcServer := grpc.NewServer()

	srv := &Server{srv: grpcServer, svc: svc}
	RegisterLoggerServer(grpcServer, srv)
	return srv
}

func (s *Server) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.srv.Serve(lis)
}

func (s *Server) FetchLogs(_ *Empty, stream Logger_FetchLogsServer) error {
	for {
		select {
		case m := <-s.svc.Logs():
			if err := stream.Send(&LogMessage{
				Label:   m.Label,
				Text:    m.Text,
				EventAt: timestamppb.New(m.EventAt),
			}); err != nil {

			}
		}
	}
}
