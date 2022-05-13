package generator

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
)

const (
	// 时间戳的长度
	timeStampSize = uint(41)
	// 高位顺序递进数的长度
	highSequenceSize = uint(8)
	// 机器编号的长度
	machineSize = uint(13)
	// 低位顺序递进数的长度
	lowSequenceSize = uint(1)

	/*
		求时间戳的最大值的公式的意思是：-1与-1乘以2的timeStampSize次方做按位异或运算

		异或运算：对比两组二进制数字的每一位上的数字，不同则在对应的结果的同一位上为1，相同则为0

		-1的二进制表示：11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111
	*/

	// 时间戳最大值
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

// Butterfly 发号器的实体类
type Butterfly struct {
	sync.Mutex
	timestamp, highSequence, machine, lowSequence int64
}

/*
	NewWithTimestamp 传入time.Now().UnixMilli()或其它int64类型的时间戳作为起始时间，获取一个发号器实例。

	请注意，该时间戳需自行持久化保存，发号器依赖于此时间戳进行后续的发号。
*/
func NewWithTimestamp(timestamp int64) (*Butterfly, error) {
	if timestamp > timestampMax {
		return nil, fmt.Errorf("timestamp[%v] can't be more than the max[%v] of timestamp", timestamp, timestampMax)
	}
	butterfly := Butterfly{
		timestamp:    timestamp,
		highSequence: 0,
		machine:      0,
		lowSequence:  0,
	}
	return &butterfly, nil
}
func NewWithId(id int64) (*Butterfly, error) {
	timestampStr := string(id)[timeStampShift:]
	timestamp, err := strconv.ParseInt(timestampStr, 0, 0)
	if err != nil {
		log.Fatalf("failed to get timestamp from the string value[%s] of timestamp: %v", timestampStr, err)
	}
	timestamp = timestamp << timeStampShift
	butterfly, _ := NewWithTimestamp(timestamp)

	return butterfly, nil
}

// NewWithTimestampAndMachineNumber 通过微秒级时间戳和机器编号构件一个发号器实例
func NewWithTimestampAndMachineNumber(timestamp, machine int64) (*Butterfly, error) {
	if machine > machineMax {
		return nil, fmt.Errorf("machine[%v] can't be more than the max[%v] of machine", machine, machineMax)
	}
	butterfly, err := NewWithTimestamp(timestamp)
	if err != nil {
		return nil, err
	}
	butterfly.machine = machine
	return butterfly, nil
}

// Generate 返回新的id给调用者
func (b *Butterfly) Generate() (int64, error) {
	b.Lock()
	// 判断低位顺序递进数是否为最大值
	if b.lowSequence == lowSequenceMax {
		// 拒绝为机器编号数值大于最大值的发号器实例继续发号
		if b.machine > machineMax {
			return 0, fmt.Errorf("the machine[%v] can't be bigger than the max[%v] of machine", b.machine, machineMax)
		}
		// 判断低位顺序递进数是否为最大值
		if b.highSequence == highSequenceMax {
			// 判断时间戳是否为最大值
			if b.timestamp == timestampMax {
				return 0, errors.New("no more id")
			} else {
				// 时间戳+1，高位顺序递进数归零
				b.timestamp++
				b.highSequence = 0
			}
		} else {
			b.highSequence++
		}
		b.lowSequence = 0
	} else {
		b.lowSequence++
	}
	// 	|是按位或运算符,当存在两个数字进行按位或运算的时候，实际进行运算的是两者的二进制数字；运算时会比较位上的数字，当两者任意一者在同一个位上存在1时，结果的该位上为1，否则为0
	id := b.timestamp<<timeStampShift | b.highSequence<<highSequenceShift | b.machine<<machineShift | b.lowSequence
	b.Unlock()
	return id, nil
}

// BatchGenerate 按数量要求批量生成符合要求数量的ID
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
