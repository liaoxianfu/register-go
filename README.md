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

