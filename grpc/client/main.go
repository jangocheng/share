package main

import (
	"flag"
	"fmt"
	"time"

	wonaming "github.com/wangming1993/share/grpc/discovery/consul"
	pb "github.com/wangming1993/share/grpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	serv = flag.String("service", "HelloService", "service name")
	reg  = flag.String("reg", "127.0.0.1:8500", "register address")
)

func main() {
	flag.Parse()
	r := wonaming.NewResolver(*serv)
	b := grpc.RoundRobin(r)

	conn, err := grpc.Dial(*reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(2 * time.Second)
	for t := range ticker.C {
		client := pb.NewHelloServiceClient(conn)
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Greeting: "world"})
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v: Reply is %s\n", t, resp.Reply)
	}
}
