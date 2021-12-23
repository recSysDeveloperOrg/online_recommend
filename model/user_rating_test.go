package model

import (
	"context"
	"testing"
)

func TestUserRatingDao_FindRatingAbove(t *testing.T) {
	dao := NewUserRatingDao()
	res, err := dao.FindRatingAbove(context.Background(), "100000000000000000158665", 3.0)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res {
		t.Logf("%+v\n", v)
	}
}
