package register_go

import "sync"

// ServiceInstancesSyncStore 用于存储Service与Instances 键值对
var ServiceInstancesSyncStore sync.Map

// ChooseFunc 负载均衡算法实现接口
type ChooseFunc func([]Instance) Instance

// Instance 注册中心中服务实例接口
type Instance interface {
	GetID() string
	SetID(ID string)
	GetServerName() string
	SetServerName(serverName string)
	GetIP() string
	SetIP(IP string)
	SetPort(port string)
	GetPort() string
	GetAddr() string
	GetMetaData() map[string]interface{}
	SetMetaData(metaData map[string]interface{})
}

// Client 注册中心客户端接口
type Client interface {
	// Close 关闭注册
	Close() error
	// RegisterServer 注册服务
	RegisterServer(info Instance) error
	// DeRegisterServer 注销服务
	DeRegisterServer(serverName string) error
	// DeRegisterInstance 注销一个实例
	DeRegisterInstance(instance Instance) error
	// GetAllInstance 获取服务名下的所有实例 首先从缓存中拿数据 如果存在就直接使用 否者再去注册中心取
	GetAllInstance(serverName string) ([]Instance, error)
	// GetOneInstance 通过chooseFunc中的规则获取一个实例 首先从缓存中拿数据 如果存在就直接使用 否者再去注册中心取
	GetOneInstance(serverName string, chooseFunc ChooseFunc) (Instance, error)
}
