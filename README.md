## register-go

利用go语言的接口特性实现服务注册与发现，可基于接口快速与不同的服务注册中心实现对接，已实现 基于`client`接口的etcd注册中心

├── README.md \
├── client.go client接口 用来定义注册中心常用的接口\
├── etcd_register etcd实现client接口和组合instance \
│ ├── etcd_client.go \
│ ├── etcd_client_test.go \
│ └── etcd_instance.go \
├── example http和grpc的简单示例\
├── go.mod \
├── go.sum \
├── grpc_resolver grpc解析服务获取instance\
│ ├── grpc_conn.go \
│ └── grpc_resolver.go \
└── instance.go 基础的instance


## 快速开始
这里以etcd作为注册中心（后期会添加更多的注册中心支持）,
### http 服务注册与发现 
  example/http/provider/provider.go
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	register "register-go"
	"register-go/etcd_register"
	"time"
)

var C register.Client
var port = "8000" // 端口

func init() {
	// 创建一个注册中心客户端
	client, err := etcd_register.NewEtcdClient([]string{"127.0.0.1:2379"}, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	// 填写本实例的相关信息
	instance := etcd_register.NewEtcdInstance("http-provider", "127.0.0.1", port, nil, 60)
    // 注册实例
	err = client.RegisterServer(instance)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 暴露客户端
	C = client
}
// http handler 
func handler(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write([]byte(port))
}
func main() {
	http.HandleFunc("/demo", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}

```

consumer 
example/http/consumer/consumer.go
```go

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	register "register-go"
	"register-go/etcd_register"
	"time"
)

var C register.Client
var port = "9000"

func init() {
	client, err := etcd_register.NewEtcdClient([]string{"127.0.0.1:2379"}, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	instance := etcd_register.NewEtcdInstance("http-consumer", "127.0.0.1", port, nil, 60)
	err = client.RegisterServer(instance)
	if err != nil {
		log.Fatalln(err)
	}
	C = client
}

func main() {
	defer func(C register.Client) {
		_ = C.Close()
	}(C)
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		// 获取一个服务实例
		instance, err := C.GetOneInstance("http-provider", nil)
		if err != nil {
			writer.WriteHeader(500)
			_, _ = writer.Write([]byte("Bad Gateway"))
			return
		}
		// get请求
		resp, err := http.Get(fmt.Sprintf("http://%s/demo", instance.GetAddr()))
		if err != nil {
			writer.WriteHeader(500)
			_, _ = writer.Write([]byte("Bad Gateway"))
			return
		}
		all, _ := ioutil.ReadAll(resp.Body)
		_, _ = writer.Write(all)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}

```

### grpc 服务注册与发现
完整代码请参见 example/grpc

proto文件
```protobuf
syntax = "proto3";

package pb;
message req{
  string msg = 1;
}
message resp{
  string info = 1;
}

service SayHello{
  rpc hello(req) returns(resp);
}
```

grpc-provider
```go
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

```

grpc-consumer

```go
package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"register-go/etcd_register"
	"register-go/example/grpc/pb"
	"register-go/grpc_resolver"
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

```