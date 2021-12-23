package model

import (
	"context"
	"testing"
)

func TestUserRatingMetaDao_FindRatingCntByUserID(t *testing.T) {
	res, err := NewUserRatingMetaDao().FindRatingCntByUserID(context.Background(), "100000000000000000158665")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
