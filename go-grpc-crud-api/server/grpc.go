package server

import (
	"context"
	"errors"
	go_grpc_crud "example/go-grpc-crud-api/proto"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GRPCServer struct {
	go_grpc_crud.UnimplementedMovieServiceServer
}

func (*GRPCServer) createMovie(DB *gorm.DB, ctx context.Context, req *go_grpc_crud.CreateMovieRequest) (*go_grpc_crud.CreateMovieResponse, error) {
	fmt.Println("Create Movie")
	movie := req.GetMovie()
	movie.Id = uuid.New().String()

	// construct for database entity
	data := Movie{
		ID:    movie.GetId(),
		Title: movie.GetTitle(),
		Genre: movie.GetGenre(),
	}

	// insert entity
	res := DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("movie creation unsuccessful")
	}

	// response -> construct grpc response
	return &go_grpc_crud.CreateMovieResponse{
		Movie: &go_grpc_crud.Movie{
			Id:    movie.GetId(),
			Title: movie.GetTitle(),
			Genre: movie.GetGenre(),
		}}, nil

}

func (*GRPCServer) GetMovie(DB *gorm.DB, ctx context.Context, req *go_grpc_crud.ReadMovieRequest) (*go_grpc_crud.ReadMovieResponse, error) {
	fmt.Println("Read Movie", req.GetId())

	var movie Movie
	res := DB.Find(&movie, "id = ?", req.GetId())
	if res.RowsAffected == 0 {
		return nil, errors.New("movie not found")
	}

	return &go_grpc_crud.ReadMovieResponse{
		Movie: &go_grpc_crud.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Genre: movie.Genre,
		},
	}, nil
}

func (*GRPCServer) GetMovies(DB *gorm.DB, ctx context.Context, req *go_grpc_crud.ReadMoviesRequest) (*go_grpc_crud.ReadMoviesResponse, error) {
	fmt.Println("Read Movies")

	var movies []*go_grpc_crud.Movie
	res := DB.Find(&movies)
	if res.RowsAffected == 0 {
		return nil, errors.New("movie not found")
	}

	return &go_grpc_crud.ReadMoviesResponse{
		Movies: movies,
	}, nil
}

func (*GRPCServer) UpdateMovie(DB *gorm.DB, ctx context.Context, req *go_grpc_crud.UpdateMovieRequest) (*go_grpc_crud.UpdateMovieResponse, error) {
	fmt.Println("Update Movie")

	var movie Movie
	reqMovie := req.GetMovie()

	res := DB.Model(&movie).
		Where("id=?", reqMovie.Id).
		Updates(Movie{Title: reqMovie.Title, Genre: reqMovie.Genre})

	if res.RowsAffected == 0 {
		return nil, errors.New("movies not found")
	}

	return &go_grpc_crud.UpdateMovieResponse{
		Movie: &go_grpc_crud.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Genre: movie.Genre,
		},
	}, nil
}

func (*GRPCServer) DeleteMovie(DB *gorm.DB, ctx context.Context, req *go_grpc_crud.DeleteMovieRequest) (*go_grpc_crud.DeleteMovieResponse, error) {
	fmt.Println("Delete Movie")

	var movie Movie
	res := DB.Where("id = ?", req.GetId()).
		Delete(&movie)
	if res.RowsAffected == 0 {
		return nil, errors.New("movie not found")
	}

	return &go_grpc_crud.DeleteMovieResponse{
		Success: true,
	}, nil
}
