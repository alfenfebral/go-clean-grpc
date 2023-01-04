package grpcdelivery

import (
	"context"
	"fmt"
	"go-clean-grpc/pkg/logger"
	proto "go-clean-grpc/todo/delivery/grpc/proto"
	models "go-clean-grpc/todo/models/http"
	todoservice "go-clean-grpc/todo/service"
	paginationutil "go-clean-grpc/utils/pagination"
)

type GRPCHandler struct {
	proto.UnimplementedTodoServer
	service todoservice.Service
}

func New(service todoservice.Service) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

func (g *GRPCHandler) Create(ctx context.Context, input *proto.TodoInput) (*proto.TodoOutput, error) {
	result, err := g.service.Create(&models.Todo{
		Title:       input.Title,
		Description: input.Description,
	})
	if err != nil {
		fmt.Println(err)
	}

	return &proto.TodoOutput{
		Id:          result.ID.Hex(),
		Title:       result.Title,
		Description: result.Description,
		CreatedAt:   result.CreatedAt.String(),
		UpdatedAt:   result.UpdatedAt.String(),
	}, nil
}

func (g *GRPCHandler) Get(ctx context.Context, input *proto.TodoGetAllInput) (*proto.TodoOutputs, error) {
	fmt.Println(input)

	page := paginationutil.CurrentPage(int(input.Page))
	perPage := paginationutil.PerPage(int(input.PerPage))
	offset := paginationutil.Offset(page, perPage)

	results, totalCount, err := g.service.GetAll(input.Q, perPage, offset)
	if err != nil {
		logger.Error(err)

		return nil, nil
	}

	pageCount := paginationutil.TotalPage(totalCount, perPage)

	var data []*proto.TodoOutput

	for _, item := range results {
		data = append(data, &proto.TodoOutput{
			Id:          item.ID.Hex(),
			Title:       item.Title,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.String(),
			UpdatedAt:   item.UpdatedAt.String(),
		})
	}

	return &proto.TodoOutputs{
		Data: data,
		Meta: &proto.Meta{
			PerPage:    int64(perPage),
			Page:       int64(page),
			PageCount:  int64(pageCount),
			TotalCount: int64(totalCount),
		},
	}, nil
}
