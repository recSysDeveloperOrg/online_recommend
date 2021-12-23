package model

import (
	"context"
	"testing"
)

func TestTagMovieDao_FindKMaxByTagID(t *testing.T) {
	dao := NewTagMovieDao()
	res, err := dao.FindKMaxByTagID(context.Background(), "61c2f1d607cb2aec6506f0f8", 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res {
		t.Logf("%+v\n", v)
	}
}
