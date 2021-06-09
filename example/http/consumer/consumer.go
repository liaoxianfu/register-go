package main

import (
	"fmt"
	register "github.com/liaoxianfu/register-go"
	"github.com/liaoxianfu/register-go/etcd_register"
	"io/ioutil"
	"log"
	"net/http"
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
