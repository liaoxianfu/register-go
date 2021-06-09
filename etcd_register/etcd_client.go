package etcd_register

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"log"
	register "register-go"
	"strings"
	"sync"
	"time"
)

// EtcdClient 抽象注册中心 register.Client 的etcd客户端实现
type EtcdClient struct {
	Endpoints   []string
	DialTimeout time.Duration
	Cli         *clientv3.Client
}

func NewEtcdClient(endpoints []string, timeout time.Duration) (*EtcdClient, error) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdClient{
		Endpoints:   endpoints,
		DialTimeout: timeout,
		Cli:         c,
	}, nil
}

func (e *EtcdClient) Close() error {
	return e.Cli.Close()
}

func (e *EtcdClient) RegisterServer(info register.Instance) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), e.DialTimeout)
	defer cancelFunc()
	instance, ok := info.(*EtcdInstance)
	if !ok {
		return errors.New("info is not a EtcdInstance")
	}
	// 服务名/ID
	k := instance.ServerName + "/" + instance.ID
	bs, err := json.Marshal(instance)
	if err != nil {
		log.Println("json.Marshal", instance, "error")
		return err
	}
	// IP:Port:metaData
	v := string(bs)
	lease, err := e.Cli.Grant(ctx, instance.TTL)
	if err != nil {
		return err
	}
	_, err = e.Cli.Put(ctx, k, v, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	keepCh, err := e.Cli.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		return err
	}
	go func(ch <-chan *clientv3.LeaseKeepAliveResponse) {
		for {
			ka := <-ch
			log.Printf("keep alive ttl = %d \n", ka.TTL)
		}
	}(keepCh)
	return nil
}

func (e *EtcdClient) DeRegisterServer(serverName string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), e.DialTimeout)
	defer cancelFunc()
	_, err := e.Cli.Delete(ctx, serverName, clientv3.WithPrefix())
	return err
}

func (e *EtcdClient) DeRegisterInstance(instance register.Instance) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), e.DialTimeout)
	defer cancelFunc()
	key := instance.GetServerName() + "/" + instance.GetID()
	_, err := e.Cli.Delete(ctx, key)
	return err
}

// WatchServer 监控sever
func (e *EtcdClient) WatchServer(serverName string) {
	log.Println("into WatchServer...")
	if _, ok := WatcherSync.Load(serverName); ok {
		log.Printf("has watched...\n")
		return
	}
	WatcherSync.Store(serverName, true)
	monitorChanel := e.Cli.Watch(context.Background(), serverName, clientv3.WithPrefix())
	for monitor := range monitorChanel {
		for _, event := range monitor.Events {
			log.Printf("type:%v,key:%v,value:%v\n", event.Type, string(event.Kv.Key), string(event.Kv.Value))
			register.ServiceInstancesSyncStore.Delete(serverName)
			_, err := e.GetAllInstance(serverName)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (e *EtcdClient) GetAllInstance(serverName string) ([]register.Instance, error) {
	load, ok := register.ServiceInstancesSyncStore.Load(serverName)
	if !ok || (ok && load == nil) {
		log.Printf("not found %s instance request register center.\n", serverName)
		ctx, cancelFunc := context.WithTimeout(context.Background(), e.DialTimeout)
		defer cancelFunc()
		response, err := e.Cli.Get(ctx, serverName, clientv3.WithPrefix())
		if err != nil {
			return nil, err
		}
		var ins []register.Instance
		if len(response.Kvs) == 0 {
			return ins, nil
		}
		for _, kv := range response.Kvs {
			instance, err := transKVToInstance(kv)
			if err != nil {
				log.Println(err)
				continue
			}
			ins = append(ins, instance)
		}
		// 进行Watch操作
		go e.WatchServer(serverName)
		register.ServiceInstancesSyncStore.Store(serverName, ins)
		return ins, nil
	}
	instances, ok := load.([]register.Instance)
	if !ok {
		return nil, errors.New("change to '[]register.Instance' error ")
	} else {
		log.Println("found server in cache")
		return instances, nil
	}
}

var mux sync.Mutex
var index = 0

func roundRobinFunc(ins []register.Instance) register.Instance {
	insLen := len(ins)
	mux.Lock()
	if index >= insLen {
		index = 0
	}
	in := ins[index]
	index++
	mux.Unlock()
	return in
}

func (e *EtcdClient) GetOneInstance(serverName string, chooseFunc register.ChooseFunc) (register.Instance, error) {
	instances, err := e.GetAllInstance(serverName)
	if err != nil {
		return nil, err
	}
	if chooseFunc == nil {
		chooseFunc = roundRobinFunc
	}
	return chooseFunc(instances), nil
}

// transKVToInstance 将kv键值对转换为 register_go.Instance
func transKVToInstance(kv *mvccpb.KeyValue) (register.Instance, error) {
	key := string(kv.Key)
	value := string(kv.Value)
	keySplit := strings.Split(key, "/")

	if len(keySplit) != 2 {
		log.Printf("key = %s not allowed", key)
		return nil, errors.New(fmt.Sprintf("key = %s not allowed", key))
	}
	var instance EtcdInstance
	err := json.Unmarshal([]byte(value), &instance)
	if err != nil {
		log.Println("json.Unmarshal ", value, " error")
		return nil, err
	}
	return &instance, nil
}
