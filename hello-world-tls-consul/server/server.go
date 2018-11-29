package main

import (
	"fmt"
	"github.com/saberuster/grpc-example/hello-world-tls-consul/internal"
	"github.com/saberuster/grpc-example/hello-world-tls-consul/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"net"
	"os"
)

const (
	address      = ":9900"
	certFilePath = "hello-world-tls-consul/public.pem"
	keyFilePath  = "hello-world-tls-consul/private.key"
)

func main() {
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)
	grpclog.Infoln("server start...")
	internal.RegisterService("hello", "localhost", 9900)
	ln, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatal(err)
		return
	}
	creds, err := credentials.NewServerTLSFromFile(certFilePath, keyFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	g := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterHelloServerServer(g, &HelloServer{})
	grpclog.Infoln(g.Serve(ln))
}

type HelloServer struct {
}

func (*HelloServer) Greeting(ctx context.Context, h *pb.Hello) (*pb.Hello, error) {
	return &pb.Hello{Message: fmt.Sprintf("reply:%s", h.Message)}, nil
}
