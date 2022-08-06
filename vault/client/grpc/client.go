package vault

import (
	. "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"grpc/vault"
	pb "grpc/vault/pb"
)

func New(conn *grpc.ClientConn) vault.Service {
	var hashEndpoint = NewClient(
		conn,
		"pb.Vault",
		"Hash",
		vault.EncodeGRPCHashReq,
		vault.DecodeGRPCHashRes,
		pb.HashRes{},
	).Endpoint()
	var validateEndpoint = NewClient(
		conn,
		"pb.Vault",
		"Validate",
		vault.EncodeGRPCValidateReq,
		vault.DecodeGRPCValidateRes,
		pb.ValidateRes{},
	).Endpoint()
	return vault.Endpoints{
		HashEndpoint:     hashEndpoint,
		ValidateEndpoint: validateEndpoint,
	}
}
