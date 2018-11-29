package main

import (
	"context"
	"fmt"
	"github.com/saberuster/grpc-example/hello-world-tls-consul/internal"
	"github.com/saberuster/grpc-example/hello-world-tls-consul/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"io/ioutil"
	"os"
)

const (
	address      = "consul:///hello"
	certFilePath = "hello-world-tls-consul/public.pem"
)

func main() {
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)
	resolver.Register(internal.NewResolveBuilder("127.0.0.1", "8500"))

	creds, err := credentials.NewClientTLSFromFile(certFilePath, "")
	if err != nil {
		log.Fatal(err)
		return
	}
	c, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}

	cc := pb.NewHelloServerClient(c)

	fmt.Println(cc.Greeting(context.TODO(), &pb.Hello{Message: "liqi"}))
}
