package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "protobuf_server/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	pb.UnimplementedCalculatorServer;
	pb.UnimplementedGreeterServiceServer;
	pb.UnimplementedFarewellServiceServer;
}

func (s *server) Add (context context.Context, req *pb.AddInterger) (*pb.AddedInteger, error) {
	return &pb.AddedInteger{
		Result: int64(req.First) + int64(req.Second),
	}, nil;
}

func (s *server) Greet (context context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Println(req.GreetRequestMessage);
	return &pb.GreetResponse{
		GreetResponseMessage: "All Good bro, you tell how what are you doing",
	}, nil;
}

func (s *server) FarewellGreetings (context context.Context, req *pb.FarewellRequestMessage) (*pb.FarewellResponseMessage, error){
	log.Println(req.Message);
	return &pb.FarewellResponseMessage{
		Message: "Thankyou very much for the farewell wishes honey",
	}, nil;
}

func main() {
	fmt.Println("")
	cert := "cert.pem"
	key := "key.pem";

	listen, err := net.Listen("tcp", "0.0.0.0:50051");
	if err != nil {
		log.Fatal("failed to listen on port", err);
	}

	credentials, err := credentials.NewServerTLSFromFile(cert, key);
	if err != nil {
		log.Fatal("failed to add cert and key usin gtl credentials", err)
	}

	log.Println("Starting grpc server")
	grpcServer := grpc.NewServer(grpc.Creds(credentials));

	pb.RegisterCalculatorServer(grpcServer, &server{});
	pb.RegisterGreeterServiceServer(grpcServer, &server{});
	pb.RegisterFarewellServiceServer(grpcServer, &server{});

	err = grpcServer.Serve(listen);
	if err != nil {
		log.Fatal("error while serving grpc server", err)
	}
}