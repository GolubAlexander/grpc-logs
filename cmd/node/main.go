package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GolubAlexander/grpc-logs/internal/server"
	"github.com/GolubAlexander/grpc-logs/internal/services/generator"
)

func main() {
	label := flag.String("label", "test", "Node's name")
	port := flag.String("port", "7001", "Node's port")
	flag.Parse()

	svc := generator.New(*label)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := svc.GenerateLogs(ctx); err != nil {
			log.Println(err)
		}
	}()

	srv := server.New(svc)
	go func() {
		if err := srv.Serve(fmt.Sprintf(":%s", *port)); err != nil {
			log.Println()
		}
	}()

	log.Println("node", *label, "started on", *port, "port")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	cancel()
	log.Println("node", *label, "stopped")
}
