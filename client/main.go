package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	pb "github.com/nnayoo/grpc-demo/proto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	addr = flag.String("addr", "localhost:50052", "the address to connect to")

	filename = flag.String("f", "", "filename")
)

const (
	timestampFormat = time.StampNano // "Jan _2 15:04:05.000"
)

func unaryCallWithMetadata(c pb.GreeterClient, message string) {
	fmt.Printf("--- unary ---\n")
	// Create metadata and context.
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Make RPC using the context with the metadata.
	var header, trailer metadata.MD
	r, err := c.Echo(ctx, &pb.Request{Name: message}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatalf("failed to call UnaryEcho: %v", err)
	}

	if t, ok := header["timestamp"]; ok {
		fmt.Printf("timestamp from header:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	} else {
		log.Fatal("timestamp expected but doesn't exist in header")
	}

	fmt.Printf("response:\n")
	fmt.Printf(" - %s\n", r.Message)

	if t, ok := trailer["timestamp"]; ok {
		fmt.Printf("timestamp from trailer:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	} else {
		log.Fatal("timestamp expected but doesn't exist in trailer")
	}
}

func main() {
	flag.Parse()
	const message = "this is examples/metadata"
	_, err := os.Stat(*filename)
	if err != nil {
		log.Fatal("file" + *filename)
	}

	data, err := ioutil.ReadFile(*filename)

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*600)
	defer cancel()
	r, err := c.Upload(ctx, &pb.Request{Name: path.Base(*filename), Data: data})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("%s", r.GetMessage())

	// r, err = c.ExCommand(ctx, &pb.Request{Name: path.Base(*filename), Data: data})
	// if err != nil {
	// 	log.Fatal("could not Ex")
	// }
	// log.Printf("%s", r.GetMessage())
	unaryCallWithMetadata(c, message)
}
