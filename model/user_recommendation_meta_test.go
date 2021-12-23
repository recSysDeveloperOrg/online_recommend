package model

import (
	"context"
	"testing"
)

const (
	testUserID = "100000099990000000158665"
)

func TestNewUserRecommendationMetaDao(t *testing.T) {
	dao := NewUserRecommendationMetaDao()
	if err := dao.AddUninterestedSet(context.Background(), testUserID,
		"100000000000000000158665", UninterestedTypeTag); err != nil {
		t.Fatal(err)
	}

	res, err := dao.FindUninterestedSet(context.Background(), testUserID, UninterestedTypeTag)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", res)

	if err := dao.AddViewLog(testUserID, "helloworld"); err != nil {
		t.Fatal(err)
	}
	viewLogs := dao.GetViewLog(testUserID)
	for _, viewLog := range viewLogs {
		t.Log(viewLog)
	}
}
