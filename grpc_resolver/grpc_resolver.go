package grpc_resolver

import (
	"google.golang.org/grpc/resolver"
	"log"
	register "register-go"
)

type GrpcResolver struct {
	target       resolver.Target
	cc           resolver.ClientConn
	addressStore map[string][]register.Instance
}
type GrpcBuilder struct {
	C register.Client
}

func (builder *GrpcBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	log.Println(target.Endpoint)
	serverName := target.Endpoint
	log.Println("target.Endpoint is ", serverName)
	// 获取所有的实例
	ins, err := builder.C.GetAllInstance(serverName)
	log.Println("register center get ", len(ins), " ins")
	if err != nil {
		return nil, err
	}
	r := &GrpcResolver{
		target: target,
		cc:     cc,
		addressStore: map[string][]register.Instance{
			target.Endpoint: ins,
		},
	}
	go r.start()
	return r, nil
}
func (r *GrpcResolver) start() {
	instances := r.addressStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(instances))
	for i, in := range instances {
		addrs[i] = resolver.Address{Addr: in.GetAddr()}
	}
	log.Println(len(addrs))
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (builder *GrpcBuilder) Scheme() string {
	return ""
}

func (r *GrpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (r *GrpcResolver) Close() {
}
