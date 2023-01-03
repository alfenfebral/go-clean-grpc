package grpcdelivery

import (
	"context"
	"fmt"
	proto "go-clean-grpc/todo/delivery/grpc/proto"
)

type GRPCHandler struct {
	proto.UnimplementedTodoServer
}

func New() *GRPCHandler {
	return &GRPCHandler{}
}

func (g *GRPCHandler) Create(ctx context.Context, input *proto.TodoInput) (*proto.TodoOutput, error) {

	fmt.Println(input.GetTitle())

	return &proto.TodoOutput{}, nil
}
