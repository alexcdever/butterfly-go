package generator

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

const (
	// the size of timestamp
	timeStampSize = uint(41)
	// the size of high sequence
	highSequenceSize = uint(8)
	// the size of machine
	machineSize = uint(13)
	// the size of low sequence
	lowSequenceSize = uint(1)

	/*
		求最大值的公式的意思是：-1与-1乘以2的Size次方做按位异或运算

		异或运算：对比两组二进制数字的每一位上的数字，不同则在对应的结果的同一位上为1，相同则为0

		-1的二进制表示：11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111
	*/

	// 时间戳最大值，将-1左移41位，则-1变成一个bit长度为41的二进制数字（左边补零）
	timestampMax = int64(-1 ^ (-1 << timeStampSize))
	// 高位顺序递进数最大值
	highSequenceMax = int64(-1 ^ (-1 << highSequenceSize))
	// 机器编号最大值
	machineMax = int64(-1 ^ (-1 << machineSize))
	// 低位顺序递进数最大值
	lowSequenceMax = int64(9)
	// 生成ID时，机器编号的数值需要左移1位
	machineShift = lowSequenceSize
	// 生成ID时，高位顺序递进数的数值需要左移14位
	highSequenceShift = machineSize + lowSequenceSize
	// 生成ID时，时间戳的数值需要左移22位
	timeStampShift = highSequenceSize + machineSize + lowSequenceSize
)

type Butterfly struct {
	sync.Mutex
	Timestamp    int64 `validate:"required,number" json:"timestamp,omitempty" yaml:"timestamp"`
	HighSequence int64 `json:"highSequence,omitempty" yaml:"highSequence"`
	Machine      int64 `validate:"required,number" json:"machine,omitempty" yaml:"machine"`
	LowSequence  int64 `json:"lowSequence,omitempty" yaml:"lowSequence"`
}

// NewFromConfigFile returns a new instance that load value from configuration file
// configFile includes name and file extension. configFile supports json and yml format
func NewFromConfigFile(configFile string) (*Butterfly, error) {
	var instance Butterfly
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config from [%s]: %v", configFile, err)
	}

	if err := viper.Unmarshal(&instance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal from [%s]: %v", configFile, err)
	}

	return &instance, nil
}

// NewWithTimestamp returns a new instance which uses the given time.Now().UnixMilli() as the timestamp
func NewWithTimestamp(timestamp int64) (*Butterfly, error) {
	if timestamp > timestampMax {
		return nil, fmt.Errorf("timestamp[%v] can't be more than the max[%v] of timestamp", timestamp, timestampMax)
	}
	butterfly, err := NewWithTimestampAndMachineNumber(timestamp, 0)
	if err != nil {
		return nil, fmt.Errorf("fialed to construct with timestamp[%v] and machine number[%v]", timestamp, 0)
	}
	return butterfly, nil
}

// NewWithNow returns a new instance that timestamp use current time
func NewWithNow() (*Butterfly, error) {
	return NewWithTimestamp(time.Now().UnixMilli())
}

// NewWithTimestampAndMachineNumber returns new instance with the given timestamp value and machine value
func NewWithTimestampAndMachineNumber(timestamp, machine int64) (*Butterfly, error) {
	if machine > machineMax {
		return nil, fmt.Errorf("machine[%v] can't be more than the max[%v] of machine", machine, machineMax)
	}
	butterfly := &Butterfly{Timestamp: timestamp, Machine: machine}
	butterfly.Machine = machine
	return butterfly, nil
}

// Generate returns new id
func (b *Butterfly) Generate() (int64, error) {
	b.Lock()
	// 判断低位顺序递进数是否为最大值
	if b.LowSequence == lowSequenceMax {
		// 拒绝为机器编号数值大于最大值的发号器实例继续发号
		if b.Machine > machineMax {
			return 0, fmt.Errorf("the machine[%v] can't be bigger than the max[%v] of machine", b.Machine, machineMax)
		}
		// 判断低位顺序递进数是否为最大值
		if b.HighSequence == highSequenceMax {
			// 判断时间戳是否为最大值
			if b.Timestamp == timestampMax {
				return 0, errors.New("no more id")
			} else {
				// 时间戳+1，高位顺序递进数归零
				b.Timestamp++
				b.HighSequence = 0
			}
		} else {
			b.HighSequence++
		}
		b.LowSequence = 0
	} else {
		b.LowSequence++
	}
	// 	|是按位或运算符,当存在两个数字进行按位或运算的时候，实际进行运算的是两者的二进制数字；运算时会比较位上的数字，当两者任意一者在同一个位上存在1时，结果的该位上为1，否则为0
	id := b.Timestamp<<timeStampShift | b.HighSequence<<highSequenceShift | b.Machine<<machineShift | b.LowSequence
	b.Unlock()
	return id, nil
}

// BatchGenerate returns the list of new id
func (b *Butterfly) BatchGenerate(count int) ([]int64, error) {
	var idList []int64
	for count > 0 {
		id, err := b.Generate()
		if err != nil {
			log.Fatalf("failed to generate id: %s", err)
			return nil, err
		}
		idList = append(idList, id)
		count--
	}
	return idList, nil
}
