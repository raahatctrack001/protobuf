package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	proto_server "proto_buf/proto/gen"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("Starting client...")

	// Load the server certificate
	certFile := "cert.pem"
	certData, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatal("could not read cert file: ", err)
	}

	// Create a cert pool with the server cert
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(certData); !ok {
		log.Fatal("failed to append cert to pool")
	}

	// Create transport credentials
	creds := credentials.NewClientTLSFromCert(certPool, "localhost")

	// Dial server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal("error while connecting client with server: ", err)
	}
	defer conn.Close()

	client := proto_server.NewCalculatorClient(conn)
	greeterClient := proto_server.NewGreeterServiceClient(conn);
	farewellClient := proto_server.NewFarewellServiceClient(conn);


	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := proto_server.AddInterger{
		First:  30,
		Second: 20,
	}

	addedInteger, err := client.Add(ctx, &req)
	if err != nil {
		log.Fatal("failed to call Add: ", err)
	}

	greetRequest := proto_server.GreetRequest {
		GreetRequestMessage : "Hanji bhai, request bhejdi hai, ab response bhejdo greeting ka",
	}
	greetingResponse, err := greeterClient.Greet(ctx, &greetRequest);
	if err != nil{
		log.Fatal("failed to get greeting message", err);
	}

	farewellRequestMessage := proto_server.FarewellRequestMessage{
		Message: "wish me for the farewell dude",
	}
	farewellResponse, err := farewellClient.FarewellGreetings(ctx, &farewellRequestMessage);
	if err != nil {
		log.Fatal("failed to load response message for farewell", err);
	}
	log.Println("greet response is here", greetingResponse);
	log.Println("Added integer result:", addedInteger)
	log.Println("response for farewell messag eis here", farewellResponse);
}
