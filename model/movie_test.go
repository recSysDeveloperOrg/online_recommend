package model

import (
	"context"
	"testing"
)

func TestMovieModelDao_FindTopKMovies(t *testing.T) {
	dao := NewMovieDao()
	movies, err := dao.FindTopKMovies(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, movie := range movies {
		t.Logf("%+v", movie)
	}
}
