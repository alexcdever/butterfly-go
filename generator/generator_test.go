package generator

import (
	"github.com/spf13/viper"
	"testing"
	"time"
)

func FuzzNewWithTimestamp(f *testing.F) {
	f.Add(time.Now().UnixMilli())
	f.Fuzz(func(t *testing.T, timestamp int64) {
		generator, err := NewWithTimestamp(timestamp)
		if err != nil || generator.Timestamp != timestamp {
			t.Errorf("failed to get instance by timestamp[%v]: %v", timestamp, err)
		}
	})
}
func TestNewWithTimestamp(t *testing.T) {
	timestamp := time.Now().UnixMilli()
	generator, err := NewWithTimestamp(timestamp)
	if err != nil || generator.Timestamp != timestamp {
		t.Errorf("failed to get instance by timestamp[%v]: %v", timestamp, err)
	}
	t.Log("successfully got instance by NewWithTimestamp")
}

func TestNewWithNow(t *testing.T) {
	gen, err := NewWithNow()
	if err != nil || gen.Timestamp > time.Now().UnixMilli() {
		t.Errorf("failed to get instance by NewWithNow: %v", err)
	}
	t.Log("successfully got instance by NewWithNow")
}

func TestNewWithTimestampAndMachineNumber(t *testing.T) {
	timestamp := time.Now().UnixMilli()
	machineId := int64(1)
	gen, err := NewWithTimestampAndMachineNumber(timestamp, machineId)
	if err != nil || gen.Timestamp > time.Now().UnixMilli() {
		t.Errorf("failed to get instance by NewWithTimestampAndMachineNumber: %v", err)
	}
	t.Log("successfully got instance by NewWithTimestampAndMachineNumber")
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

func TestReadConfig(t *testing.T) {
	var jsonConfig Butterfly
	var ymlConfig Butterfly
	jsonViper := viper.New()
	ymlViper := viper.New()
	jsonViper.SetConfigFile("test.json")
	ymlViper.SetConfigFile("test.yml")
	if err := jsonViper.ReadInConfig(); err != nil {
		t.Errorf("failed to read config from json: %v", err)
	}
	if err := ymlViper.ReadInConfig(); err != nil {
		t.Errorf("failed to read config from yml: %v", err)
	}

	if err := jsonViper.Unmarshal(&jsonConfig); err != nil {
		t.Errorf("failed to unmarshal from json: %v", err)
	}
	if err := ymlViper.Unmarshal(&ymlConfig); err != nil {
		t.Errorf("failed to unmarshal from json: %v", err)
	}

	jsonId, err := jsonConfig.Generate()
	if err != nil {
		t.Errorf("the generator form json failed to generate new id: %v", err)
	}
	ymlId, err := ymlConfig.Generate()
	if err != nil {
		t.Errorf("the generator form yml failed to generate new id: %v", err)
	}
	if jsonId == ymlId {
		t.Log("generator can load config from file")
		t.Logf("json config: %p", &jsonConfig)
		t.Logf("yml config: %p", &ymlConfig)
	} else {
		t.Fatalf("json config[%v] is different with yml config[%v]", &jsonConfig, &ymlConfig)
	}
}
