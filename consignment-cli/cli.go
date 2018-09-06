//consignment-cli/main.go

package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/midnightrun/shippy/consignment-service/proto/consignment"
	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

// parseFile parses a filepath and gives the file back to the caller
func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment

	data, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewShippingServiceClient(conn)

	file := defaultFilename

	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)

	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}

	log.Printf("Created: %t", r.Created)

	consignments, err := client.GetConsignment(context.Background(), &pb.GetRequest{})

	if err != nil {
		log.Fatalf("Could not get consignments: %v", err)
	}

	for _, v := range consignments.Consignments {
		log.Printf("Consignment: %v", v)
	}
}
