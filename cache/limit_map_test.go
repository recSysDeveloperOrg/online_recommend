package cache

import "testing"

const (
	opTypeWrite = "write"
	opTypeRead  = "read"
)

type cases struct {
	opType     string
	key        string
	value      string
	shouldFind bool
}

func TestLimitMap_Add(t *testing.T) {
	limitSize := 4
	lm := NewLimitMap(limitSize)
	testCases := []*cases{
		{opType: opTypeWrite, key: "hello", value: "world"},
		{opType: opTypeWrite, key: "hello", value: "u"},
		{opType: opTypeWrite, key: "hello", value: "fuckyou"},
		{opType: opTypeWrite, key: "hello", value: "ha"},
		{opType: opTypeRead, key: "hello", value: "world", shouldFind: true},
		{opType: opTypeRead, key: "hello", value: "u", shouldFind: true},
		{opType: opTypeRead, key: "hello", value: "fuckyou", shouldFind: true},
		{opType: opTypeRead, key: "hello", value: "ha", shouldFind: true},
		{opType: opTypeWrite, key: "hello", value: "123"},
		{opType: opTypeRead, key: "hello", value: "world", shouldFind: false},
		{opType: opTypeRead, key: "?", shouldFind: false},
	}
	for _, testCase := range testCases {
		if testCase.opType == opTypeWrite {
			lm.Add(testCase.key, testCase.value)
		} else {
			if testCase.shouldFind != lm.CheckExist(testCase.key, testCase.value) {
				t.Fatalf("value %s failed\n", testCase.value)
			}
		}
	}
}
