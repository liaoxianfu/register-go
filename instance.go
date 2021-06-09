package register_go

// BaseInstance 是 register_go.Instance 的基础实现
type BaseInstance struct {
	ID         string
	ServerName string
	IP         string
	Port       string
	MetaData   map[string]interface{}
}

func (b *BaseInstance) GetAddr() string {
	return b.IP + ":" + b.Port
}

func (b *BaseInstance) SetPort(port string) {
	b.Port = port
}

func (b *BaseInstance) GetPort() string {
	return b.Port
}

func (b *BaseInstance) GetID() string {
	return b.ID
}

func (b *BaseInstance) SetID(ID string) {
	b.ID = ID
}

func (b *BaseInstance) GetServerName() string {
	return b.ServerName
}

func (b *BaseInstance) SetServerName(serverName string) {
	b.ServerName = serverName
}

func (b *BaseInstance) GetIP() string {
	return b.IP
}

func (b *BaseInstance) SetIP(IP string) {
	b.IP = IP
}

func (b BaseInstance) GetMetaData() map[string]interface{} {
	return b.MetaData
}

func (b BaseInstance) SetMetaData(metaData map[string]interface{}) {
	b.MetaData = metaData
}
