package main

import (
	"database/sql"
	"log"
	"net"

	proto "server/user"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedUserServiceServer
	db *sql.DB
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, NewServer())

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
