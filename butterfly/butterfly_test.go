package butterfly

import (
	"testing"
	"time"
)

func FuzzNewWithTimestamp(f *testing.F) {
	f.Add(time.Now().UnixMilli())
	f.Fuzz(func(t *testing.T, timestamp int64) {
		generator := NewWithTimestamp(timestamp)
		if generator.Timestamp != timestamp {
			t.Errorf("failed to get instance by timestamp[%v]", timestamp)
		}
	})
}
func TestNewWithTimestamp(t *testing.T) {
	timestamp := time.Now().UnixMilli()
	generator := NewWithTimestamp(timestamp)
	if generator.Timestamp != timestamp {
		t.Errorf("failed to get instance by timestamp[%v]", timestamp)
	}
	t.Log("successfully got instance by NewWithTimestamp")
}

func TestNewWithNow(t *testing.T) {
	gen := NewWithNow()
	if gen.Timestamp > time.Now().UnixMilli() {
		t.Errorf("failed to get instance by NewWithNow")
	}
	t.Log("successfully got instance by NewWithNow")
}

func TestNewWithTimestampAndMachineNumber(t *testing.T) {
	timestamp := time.Now().UnixMilli()
	machineId := int64(1)
	gen := NewWithTimestampAndMachineNumber(timestamp, machineId)
	if gen.Timestamp > time.Now().UnixMilli() {
		t.Errorf("failed to get instance by NewWithTimestampAndMachineNumber")
	}
	t.Log("successfully got instance by NewWithTimestampAndMachineNumber")
}
func TestButterfly_Generate(t *testing.T) {
	count := 20
	generator := NewWithTimestamp(time.Now().UnixMilli())

	var result map[int64]interface{}
	result = make(map[int64]interface{})
	for i := 0; i < count; i++ {
		id := generator.Generate()
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
	generator := NewWithTimestamp(time.Now().UnixMilli())

	var result map[int64]interface{}
	result = make(map[int64]interface{})
	idList := generator.BatchGenerate(count)

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

// execute `go test -bench=. -benchmem` in the folder of where does this file exists
func Benchmark_GenerateWithoutDB(b *testing.B) {

	millionCount := b.N
	b.ResetTimer()
	generator := NewWithTimestamp(time.Now().UnixMilli())

	var result map[int64]interface{}
	result = make(map[int64]interface{})
	for i := 0; i < millionCount; i++ {
		id := generator.Generate()
		result[id] = 0
		b.Log(id)
	}
	if len(result) != millionCount {
		b.Errorf("the count of id is not correct, expected [%v], but [%v] ", millionCount, len(result))
	} else {
		b.Log("test successfully")
	}
}

func BenchmarkButterfly_GenerateWithDB(b *testing.B) {

}

func TestNewFromConfigFile(t *testing.T) {
	jsonConfig, err := NewFromConfigFile("config.json")
	ymlConfig, err := NewFromConfigFile("config.yml")

	jsonId := jsonConfig.Generate()
	if err != nil {
		t.Errorf("the butterfly form json failed to generate new id: %v", err)
	}
	ymlId := ymlConfig.Generate()
	if err != nil {
		t.Errorf("the butterfly form yml failed to generate new id: %v", err)
	}
	if jsonId == ymlId {
		t.Log("butterfly can load config from file")
		t.Logf("json config: %p", &jsonConfig)
		t.Logf("yml config: %p", &ymlConfig)
	} else {
		t.Fatalf("json config[%v] is different with yml config[%v]", &jsonConfig, &ymlConfig)
	}
}
