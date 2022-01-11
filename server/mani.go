/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"path"
	"time"

	pb "github.com/nnayoo/grpc-demo/proto/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

var (
	port   = flag.Int("port", 50052, "The server port")
	logger *log.Logger
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

const (
	timestampFormat = time.StampNano
)

// Upload implements helloworld.GreeterServer
func (s *server) Upload(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	log.Printf("Received: %v", in.GetName())
	err := ioutil.WriteFile(in.GetName(), in.GetData(), 0655)
	if err != nil {
		return &pb.Reply{Message: "send " + in.GetName() + " failed"}, nil
	} else {
		return &pb.Reply{Message: "send " + in.GetName() + " sucess"}, nil
	}
}

func (s *server) ExCommod(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	if path.Ext(in.Name) == ".sh" {
		command := exec.Command("/bin/sh", "-c", in.GetName())
		stdout, _ := command.StdoutPipe()
		stderr, _ := command.StderrPipe()
		err := command.Start()
		if err != nil {
			log.Printf(err.Error())
		}
		out_st, _ := ioutil.ReadAll(stdout)
		out_err, _ := ioutil.ReadAll(stderr)
		stdout.Close()

		var result string = "sucess"
		if string(out_st) != "" {
			logger.Println("result " + string(out_st))

		}
		if string(out_err) != "" {
			logger.Println("error: " + string(out_err))
			result = "failed"

		} else {
			logger.Println("sucess")

		}

		return &pb.Reply{Message: result + " " + in.GetName()}, nil
	} else {
		logger.Println("invalid command or not script file")
		return &pb.Reply{Message: "invalid command or not script file " + in.GetName()}, nil
	}

}

func (s *server) Echo(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		grpc.SetTrailer(ctx, trailer)
	}()

	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "UnaryEcho: failed to get metadata")
	}
	if t, ok := md["timestamp"]; ok {
		fmt.Printf("timestamp from metadata:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}

	// Create and send header.
	header := metadata.New(map[string]string{"location": "MTV", "timestamp": time.Now().Format(timestampFormat)})
	grpc.SendHeader(ctx, header)

	fmt.Printf("request received: %v, sending echo\n", in)

	return &pb.Reply{Message: in.Name}, nil
}

func main() {

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var option = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(209715200),
		grpc.MaxSendMsgSize(209715200),
	}

	s := grpc.NewServer(option...)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
