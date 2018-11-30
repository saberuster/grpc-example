package main

import (
	"flag"
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
	"strconv"
)

var (
	address  = flag.String("address", "127.0.0.1", "listener address")
	port     = flag.Int("port", 9090, "port")
	certPath = flag.String("crt", "public.pem", "crt")
	keyPath  = flag.String("key", "private.key", "key")
)

func main() {
	flag.Parse()
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)
	grpclog.Infoln("server start...")
	internal.RegisterService("hello", *address, *port)
	ln, err := net.Listen("tcp4", net.JoinHostPort(*address, strconv.Itoa(*port)))
	if err != nil {
		log.Fatal(err)
		return
	}
	creds, err := credentials.NewServerTLSFromFile(*certPath, *keyPath)
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
