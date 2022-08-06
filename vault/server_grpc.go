package vault

import (
	"context"
	. "github.com/go-kit/kit/transport/grpc"
	pb "grpc/vault/pb"
)

type grpcServer struct {
	hash     Handler
	validate Handler
}

func (s *grpcServer) Hash(ctx context.Context, r *pb.HashReq) (*pb.HashRes, error) {
	_, resp, err := s.hash.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.HashRes), nil
}

func (s *grpcServer) Validate(ctx context.Context, r *pb.ValidateReq) (*pb.ValidateRes, error) {
	_, resp, err := s.validate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ValidateRes), nil
}

func EncodeGRPCHashReq(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(hashReq)
	return &pb.HashReq{Password: req.Password}, nil
}

func DecodeGRPCHashReq(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.HashReq)
	return hashReq{req.Password}, nil
}

func EncodeGRPCHashRes(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(hashRes)
	return &pb.HashRes{Hash: res.Hash, Err: res.Err}, nil
}

func DecodeGRPCHashRes(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.HashRes)
	return hashRes{res.Hash, res.Err}, nil
}

func EncodeGRPCValidateReq(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(validateReq)
	return &pb.ValidateReq{Password: req.Password, Hash: req.Hash}, nil
}

func DecodeGRPCValidateReq(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ValidateReq)
	return validateReq{Password: req.Password, Hash: req.Hash}, nil
}

func EncodeGRPCValidateRes(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(validateRes)
	return &pb.ValidateRes{Valid: res.Valid}, nil
}

func DecodeGRPCValidateRes(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.ValidateRes)
	return validateRes{Valid: res.Valid}, nil
}

func NewGRPCServer(ctx context.Context, endpoints Endpoints) pb.VaultServer {
	return &grpcServer{
		hash: NewServer(
			endpoints.HashEndpoint,
			DecodeGRPCHashReq,
			EncodeGRPCHashRes,
		),
		validate: NewServer(
			endpoints.ValidateEndpoint,
			DecodeGRPCValidateReq,
			EncodeGRPCValidateRes,
		),
	}
}
