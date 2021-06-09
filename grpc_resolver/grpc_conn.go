package grpc_resolver

import "google.golang.org/grpc"

var ConnMap map[string]*grpc.ClientConn
var DefaultGRPCConn = &GRPCConn{}

type IGRPCConn interface {
	GetGRPCConn(target string, ops ...grpc.DialOption) (*grpc.ClientConn, error)
}

type GRPCConn struct {
}

func (conn *GRPCConn) GetGRPCConn(serviceName string, ops ...grpc.DialOption) (*grpc.ClientConn, error) {
	if ops == nil || len(ops) == 0 {
		ops = []grpc.DialOption{
			grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy":"round_robin"}`),
			grpc.WithInsecure(),
		}
	}
	if ConnMap[serviceName] == nil {
		return grpc.Dial(serviceName, ops...)
	}
	return ConnMap[serviceName], nil
}
