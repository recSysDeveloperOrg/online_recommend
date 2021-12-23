package model

import (
	"context"
	"testing"
)

func TestTagUserDao_FindKMaxUserTags(t *testing.T) {
	dao := NewTagUserDao()
	res, err := dao.FindKMaxUserTags(context.Background(), "100000000000000000158665", 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res {
		t.Logf("%+v\n", v)
	}
}
