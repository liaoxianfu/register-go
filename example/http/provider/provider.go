package main

import (
	"fmt"
	register "github.com/liaoxianfu/register-go"
	"github.com/liaoxianfu/register-go/etcd_register"
	"log"
	"net/http"
	"time"
)

var C register.Client
var port = "8000"

func init() {
	client, err := etcd_register.NewEtcdClient([]string{"127.0.0.1:2379"}, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	instance := etcd_register.NewEtcdInstance("http-provider", "127.0.0.1", port, nil, 60)
	err = client.RegisterServer(instance)
	if err != nil {
		log.Fatalln(err.Error())
	}
	C = client
}
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
