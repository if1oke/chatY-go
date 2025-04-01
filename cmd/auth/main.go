package main

import (
	"chatY-go/internal/api/grpc/auth/proto"
	"chatY-go/internal/application/auth"
	"chatY-go/pkg/logger"
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	var app auth.IApplication
	var logger = logger.NewLogger()

	app = auth.NewApplication(logger)
	app.Init()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", app.Config().ServerAddress(), app.Config().AuthPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	proto.RegisterAuthServiceServer(app.GRPCServer(), app.AuthHandler())

	logger.Infof("Auth gRPC server listening on %s:%s", app.Config().ServerAddress(), app.Config().AuthPort())
	reflection.Register(app.GRPCServer())
	err = app.GRPCServer().Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
