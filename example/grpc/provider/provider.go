package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"register-go/etcd_register"
	"register-go/example/grpc/pb"
	"time"
)

var port = "5000"

type HelloStruct struct {
}

func (s *HelloStruct) Hello(ctx context.Context, req *pb.Req) (*pb.Resp, error) {
	msg := req.Msg + port
	return &pb.Resp{
		Info: msg,
	}, nil
}

func init() {
	client, err := etcd_register.NewEtcdClient([]string{"127.0.0.1:2379"}, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	instance := etcd_register.NewEtcdInstance("grpc-provider", "127.0.0.1", port, nil, 10)
	err = client.RegisterServer(instance)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatalln(err)
	}
	server := grpc.NewServer()
	pb.RegisterSayHelloServer(server, &HelloStruct{})
	err = server.Serve(listen)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
