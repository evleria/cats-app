// Package grpc encapsulates work with gRPC
package grpc

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/evleria/mongo-crud/internal/pb"
	"github.com/evleria/mongo-crud/internal/repository"
	"github.com/evleria/mongo-crud/internal/repository/entities"
	"github.com/evleria/mongo-crud/internal/service"
)

// CatsService grpc service implementation of pb.CatsServiceServer
type CatsService struct {
	pb.UnimplementedCatsServiceServer
	service service.Cats
}

// NewCatsService returns a new pb.CatsService
func NewCatsService(catsService service.Cats) pb.CatsServiceServer {
	return &CatsService{
		service: catsService,
	}
}

// GetAllCats fetches all cats from cats collection
func (s *CatsService) GetAllCats(ctx context.Context, _ *empty.Empty) (*pb.GetAllResponse, error) {
	cats, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GetAllResponse{
		Cats: mapCats(cats),
	}
	return response, nil
}

// GetCat fetches a cat from cats collection by ID
func (s *CatsService) GetCat(ctx context.Context, request *pb.GetCatRequest) (*pb.GetCatResponse, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	cat, err := s.service.GetOne(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GetCatResponse{
		Cat: mapCat(cat),
	}
	return response, nil
}

// AddNewCat creates new cat in cats collection
func (s *CatsService) AddNewCat(ctx context.Context, request *pb.AddNewCatRequest) (*pb.AddNewCatResponse, error) {
	id, err := s.service.CreateNew(ctx, request.Name, request.Color, int(request.Age), request.Price)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := &pb.AddNewCatResponse{
		Id: id.String(),
	}
	return response, nil
}

// DeleteCat removes cat from cats collection by ID
func (s *CatsService) DeleteCat(ctx context.Context, request *pb.DeleteCatRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = s.service.Delete(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &empty.Empty{}, nil
}

// UpdatePrice updates price of a cat by id
func (s *CatsService) UpdatePrice(ctx context.Context, request *pb.UpdatePriceRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = s.service.UpdatePrice(ctx, id, request.Price)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func mapCat(cat entities.Cat) *pb.Cat {
	return &pb.Cat{
		Id:    cat.ID.String(),
		Name:  cat.Name,
		Color: cat.Color,
		Age:   int64(cat.Age),
		Price: cat.Price,
	}
}

func mapCats(cats []entities.Cat) []*pb.Cat {
	result := make([]*pb.Cat, 0, len(cats))
	for _, cat := range cats {
		result = append(result, mapCat(cat))
	}
	return result
}
