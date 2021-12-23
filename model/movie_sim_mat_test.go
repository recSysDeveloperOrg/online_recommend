package model

import (
	"context"
	"log"
	"testing"
)

func TestMovieSimMatDao_FindByMovieIDs(t *testing.T) {
	dao := NewMovieSimMatDao()
	topMovies, err := NewMovieDao().FindTopKMovies(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	res, err := dao.FindByMovieIDs(context.Background(), topMovies)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res {
		for _, w := range v.To {
			log.Printf("%+v\n", w)
		}
	}
}
