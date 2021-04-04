package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GolubAlexander/grpc-logs/internal/services/collector"
)

func main() {

	addrs := os.Args[1:]

	logger := log.New(os.Stderr, "", log.LstdFlags)

	svc := collector.New(logger)

	for _, addr := range addrs {
		if err := svc.ConnectNode(addr, fmt.Sprintf("node %s", addr)); err != nil {
			logger.Fatal("fail to connect to node 7001")
		}
	}

	svc.FetchLogs()
}
