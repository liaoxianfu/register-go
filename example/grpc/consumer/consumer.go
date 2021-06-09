package main

import (
	"context"
	"github.com/liaoxianfu/register-go/etcd_register"
	"github.com/liaoxianfu/register-go/example/grpc/pb"
	"github.com/liaoxianfu/register-go/grpc_resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"time"
)

var port = "7000"

func init() {
	client, err := etcd_register.NewEtcdClient([]string{"127.0.0.1:2379"}, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	instance := etcd_register.NewEtcdInstance("grpc-consumer", "127.0.0.1", port, nil, 10)
	err = client.RegisterServer(instance)
	if err != nil {
		log.Fatalln(err)
	}
	builder := grpc_resolver.GrpcBuilder{C: client}
	resolver.Register(&builder)
}

func main() {
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		// ops 默认的参数如下
		ops := []grpc.DialOption{
			// 使用轮询
			grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy":"round_robin"}`),
			// 不安全的连接
			grpc.WithInsecure(),
		}
		// 获取连接
		conn, err := grpc_resolver.DefaultGRPCConn.GetGRPCConn("grpc-provider", ops...)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadGateway)
			_, _ = writer.Write([]byte("Bad GateWay"))
			return
		}
		sayHelloClient := pb.NewSayHelloClient(conn)
		// 使用超时时间为2s
		deadline, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
		defer cancelFunc()
		hello, err := sayHelloClient.Hello(deadline, &pb.Req{Msg: "hello"})
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadGateway)
			_, _ = writer.Write([]byte("Bad GateWay"))
			return
		}
		_, _ = writer.Write([]byte(hello.Info))
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalln(err)
	}

}
