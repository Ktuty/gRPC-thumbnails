package client

import (
	"google.golang.org/grpc"
)

func NewClient(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
}
