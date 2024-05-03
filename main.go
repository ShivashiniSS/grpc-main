package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"grpc/student-api/database"
	"grpc/student-api/pb"
	"grpc/student-api/studentdex"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("failed to create listener:", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pool, err := database.Connect(context.Background())
	if err != nil {
		log.Fatalln("failed to connect to database:", err)
	}

	repository := studentdex.NewRepository(pool)
	server := studentdex.NewServer(repository)

	pb.RegisterStudentdexServer(s, server)
	if err := s.Serve(listener); err != nil {
		log.Fatalln("failed to serve:", err)
	}
}
