package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/grpc-go"
	pb "github.com/nnayoo/grpc-demo/proto/proto"
)

var (
	addr = flag.String("addr", "localhost:50052", "the address to connect to")

	filename = flag.String("f", "", "filename")
)

func main() {
	_, err := os.Stat(*filename)
	if err != nil {
		log.Fatal("file not found")
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

	r, err = c.ExCommand(ctx, &pb.Request{Name: path.Base(*filename), Data: data})
	if err != nil {
		log.Fatal("could not Ex")
	}
	log.Printf("%s", r.GetMessage())

}
