package generator

import (
	"testing"
	"time"
)

func TestNewWithTimestamp(t *testing.T) {
	timestamp := time.Now().UnixMilli()
	generator, err := NewWithTimestamp(timestamp)
	if err != nil {
		t.Errorf("failed to construct generator: %v", err)
	}
	if generator.timestamp == timestamp {
		t.Log("test successfully")
	} else {
		t.Errorf("timestamp is different: origin[%d]:generator[%d]", timestamp, generator.timestamp)
	}
}

func TestNewWithId(t *testing.T) {
	generator, err := NewWithTimestamp(time.Now().UnixMilli())
	if err != nil {
		t.Errorf("failed to construct generator: %v", err)
	}
	id, err := generator.Generate()
	if err != nil {
		t.Errorf("failed to generate id: %v", err)
	}
	gen2, err := NewWithId(id)
	if err != nil {
		t.Errorf("failed to construct the second generator: %v", err)
	}
	id2, err := gen2.Generate()
	if err != nil {
		t.Errorf("failed to generate id by gen2: %v", err)
	}
	if id2 > id {
		t.Log("test successfully")
	} else {
		t.Errorf("the id2[%d] must be bigger than id[%d]", id2, id)
	}

}

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
