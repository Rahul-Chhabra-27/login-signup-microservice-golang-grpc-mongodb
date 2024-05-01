package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	userproto "rahulchhabra.io/proto/user"
)

func createuser(client userproto.UserServiceClient) {
	// Call the CreateUser function on the server and pass the request to it using the client we created above.
	response, err := client.CreateUser(context.Background(), &userproto.CreateUserRequest{
		User: &userproto.User{
			Firstname: "alex",
			Lastname:  "Doe",
			Email:     "alenxc@mail.com",
			Username:  "rahulc",
		},
		Password: "password",
	})
	//Check for errors
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}
	// Print the response
	fmt.Println("User created: ", response.GetUser().GetId())
}

func loginuser(client userproto.UserServiceClient) {
	// Call the AuthenticateUser function on the server and pass the request to it using the client we created above.
	response, err := client.AuthenticateUser(context.Background(), &userproto.AuthenticateUserRequest{
		Email: "rahul.c@prograd.org",
		Password: "password",
	})
	//Check for errors
	if err != nil {
		log.Fatalf("Could not authenticate user: %v", err)
	}
	// Print the response
	fmt.Println("User authenticated: \n", response.Message);
	fmt.Println("Token Generated :- ", response.AuthToken);
}
func main() {
	// Create a connection to the server
	// grpc.Dial is a function that creates a connection to the server using the gRPC protocol
	// grpc.Dial takes the address of the server and the credentials
	// In this case, we are using insecure credentials
	// This is because we are not using TLS((Transport Layer Security) certificates)
	connection, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to the server: %v", err)
	}
	// Close the connection when the function exits
	defer connection.Close()

	// Create a new client
	client := userproto.NewUserServiceClient(connection)

	//createuser(client)
	loginuser(client)
}
