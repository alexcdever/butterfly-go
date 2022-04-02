package generator

import (
	"testing"
	"time"
)

func TestButterfly_Generate(t *testing.T) {
	count := 20
	generator, err := NewWithTimestamp(time.Now().UnixMilli())
	if err != nil {
		t.Errorf("failed to get generator: %s", err)
	}
	var result map[int64]interface{}
	result = make(map[int64]interface{})
	for i := 0; i < count; i++ {
		id, _ := generator.Generate()
		result[id] = 0
		t.Log(id)
	}
	if len(result) != count {
		t.Errorf("the count of id is not correct, expected [%v], but [%v] ", count, len(result))
	} else {
		t.Log("test successfully")
	}
}

func TestButterfly_BatchGenerate(t *testing.T) {
	count := 20
	generator, err := NewWithTimestamp(time.Now().UnixMilli())
	if err != nil {
		t.Errorf("failed to get generator: %s", err)
	}
	var result map[int64]interface{}
	result = make(map[int64]interface{})
	idList, _ := generator.BatchGenerate(count)

	for _, v := range idList {
		result[v] = 0
		t.Log(v)
	}

	if len(result) != count {
		t.Errorf("the size of id list is not correct, expected [%v], but [%v] ", count, len(result))
	} else {
		t.Log("test successfully")
	}
}
