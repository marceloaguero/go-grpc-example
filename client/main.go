package main

import (
	"io"
	"log"

	pb "github.com/marceloaguero/go-grpc-example/customer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

// createCustomer calls the RPC method CreateCustomer of CustomerServer
func createCustomer(client pb.CustomerClient, customer *pb.CustomerRequest) {
	resp, err := client.CreateCustomer(context.Background(), customer)
	if err != nil {
		log.Fatalf("could not create Customer: %v", err)
	}
	if resp.Success {
		log.Printf("A new customer had been added with id: %d", resp.Id)
	}
}

// getCustomers calls the RPC method GetCustomers of CustomerServer
func getCustomers(client pb.CustomerClient, filter *pb.CustomerFilter) {
	// calling the streaming API
	stream, err := client.GetCustomers(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}
		log.Printf("Customer: %v", customer)
	}
}

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Creates a new CustomerClient
	client := pb.NewCustomerClient(conn)

	customer := &pb.CustomerRequest{
		Id:    101,
		Name:  "Marcelo Aguero",
		Email: "marceloaguero@gmail.com",
		Phone: "+5492604680717",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "Av. Velez Sarsfield 1095",
				City:              "San Rafael",
				State:             "Mendoza",
				Zip:               "5600",
				IsShippingAddress: true,
			},
			&pb.CustomerRequest_Address{
				Street:            "Horacio Morales esq. Ej. de los Andes",
				City:              "Rama Caida",
				State:             "Mendoza",
				Zip:               "5600",
				IsShippingAddress: false,
			},
		},
	}

	// Create a new customer
	createCustomer(client, customer)

	customer = &pb.CustomerRequest{
		Id:    102,
		Name:  "Aurelia Giannattasio",
		Email: "aurelia@demetersrl.com",
		Phone: "",
		Addresses: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "Av. Velez Sarsfield 1095",
				City:              "San Rafael",
				State:             "Mendoza",
				Zip:               "5600",
				IsShippingAddress: true,
			},
		},
	}

	// Create a new customer
	createCustomer(client, customer)
	// Filter with an empty Keyword
	filter := &pb.CustomerFilter{Keyword: ""}
	getCustomers(client, filter)
}
