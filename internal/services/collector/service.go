package collector

import (
	"context"
	"errors"
	"fmt"
	"io"

	"google.golang.org/grpc"

	pb "github.com/GolubAlexander/grpc-logs/internal/server"
)

type Service struct {
	// TODO: mutex
	nodes map[string]pb.LoggerClient
	log   Logger
}

type Logger interface {
	Printf(format string, v ...interface{})
}

func New(log Logger) *Service {
	return &Service{
		nodes: make(map[string]pb.LoggerClient),
		log:   log,
	}
}

func (svc *Service) ConnectNode(addr string, label string) error {
	if svc.nodeExists(label) {
		return fmt.Errorf("node %s is already existed", label)
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewLoggerClient(conn)
	svc.nodes[label] = client
	return nil
}

func (svc *Service) FetchLogs() {
	for label, node := range svc.nodes {
		stream, err := node.FetchLogs(context.Background(), &pb.Empty{})
		if err != nil {
			svc.log.Printf("label=%s; error=%s;\n", label, err)
			continue
		}
		go func(label string) {
			for {
				m, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					svc.log.Printf("label=%s; error=%s;\n", label, err)
				}
				svc.log.Printf("label=%s; eventAt=%s; text=%s\n", m.GetLabel(), m.GetEventAt().AsTime(), m.GetText())
			}
		}(label)
	}
	select {}
}

func (svc *Service) nodeExists(label string) bool {
	_, found := svc.nodes[label]
	return found
}
