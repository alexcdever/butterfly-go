package main

import (
	"testing"
	"time"
)

func TestNext(t *testing.T) {
	count := 20
	generator := GetGenerator(time.Now().UnixNano())
	var result map[int64]interface{}
	result = make(map[int64]interface{})
	for i := 0; i < count; i++ {
		id, _ := generator.Next()
		result[id] = 0
		t.Log(id)
	}
	if len(result) != count {
		t.Errorf("the count of id is not correct, expected [%v], but [%v] ", count, len(result))
	} else {
		t.Log("test successfully")
	}
}
