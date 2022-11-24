package rpc

import (
	"inetmock.icb4dc0.de/inetmock/netflow"
	rpcv1 "inetmock.icb4dc0.de/inetmock/pkg/rpc/v1"
)

// var _ rpcv1.NetFlowControlServiceServer = (*netMonServer)(nil)

func NewNetFlowControlServiceServer(fw *netflow.Firewall, nat *netflow.NAT) rpcv1.NetFlowControlServiceServer {
	return &netMonServer{
		fw:  fw,
		nat: nat,
	}
}

type netMonServer struct {
	rpcv1.UnimplementedNetFlowControlServiceServer
	fw  *netflow.Firewall
	nat *netflow.NAT
}
