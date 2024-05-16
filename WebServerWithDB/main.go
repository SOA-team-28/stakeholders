package main

import (
	"database-example/db"
	user_service "database-example/proto/user"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startServer() {
	database := db.InitDB()
	if database == nil {
		log.Fatal("FAILED TO CONNECT TO DB")
	}
	lis, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)

	// Bootstrap gRPC service server and respond to request.

	user_service.RegisterUserServiceServer(grpcServer, Server{})
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)

	router := mux.NewRouter().StrictSlash(true)

	//dodati handlere

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8086", router))
}

type Server struct {
	user_service.UnimplementedUserServiceServer
}

func main() {

	startServer()
}
