package etcd_register

import (
	"github.com/google/uuid"
	register "github.com/liaoxianfu/register-go"
	"sync"
)

// EtcdInstance etcd具体实例
type EtcdInstance struct {
	register.BaseInstance
	TTL int64
}

// WatcherSync 使用并发安全map
var WatcherSync sync.Map

func (ins *EtcdInstance) SetTTL(ttl int64) {
	ins.TTL = ttl
}
func (ins *EtcdInstance) GetTTL() int64 {
	return ins.TTL
}

func NewEtcdInstance(serverName, ip, port string, metaData map[string]interface{}, ttl int64) register.Instance {
	return &EtcdInstance{
		BaseInstance: register.BaseInstance{
			ID:         uuid.NewString(),
			ServerName: serverName,
			IP:         ip,
			Port:       port,
			MetaData:   metaData,
		},
		TTL: ttl,
	}
}
