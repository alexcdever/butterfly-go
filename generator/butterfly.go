package generator

import (
	"errors"
	"sync"
)

const (
	timeStampSize    = uint(41)
	highSequenceSize = uint(8)
	machineSize      = uint(13)
	lowSequenceSize  = uint(1)
	/*
		等价于：-1与-1乘以2的timeStampSize次方做按位异或运算

		异或运算：对比两组二进制数字的每一位上的数字，不同则在对应的结果的同一位上为1，相同则为0

		-1的二进制表示：11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111
	*/
	timestampMax      = int64(-1 ^ (-1 << timeStampSize))
	highSequenceMax   = int64(-1 ^ (-1 << highSequenceSize))
	machineMax        = int64(-1 ^ (-1 << machineSize))
	lowSequenceMax    = int64(9)
	machineShift      = lowSequenceSize
	highSequenceShift = machineSize + lowSequenceSize
	timeStampShift    = highSequenceSize + machineSize + lowSequenceSize
)

type Butterfly struct {
	sync.Mutex
	timeStamp    int64
	highSequence int64
	machine      int64
	lowSequence  int64
}

/*
	传入time.Now().UnixNano()或其它int64类型的时间戳数字，获取一个发号器实例。

	请注意，该时间戳请自行持久化保存，发号器依赖于此时间戳进行发号。
*/
func New(timeStamp int64) *Butterfly {
	butterfly := Butterfly{
		timeStamp:    timeStamp,
		highSequence: 0,
		machine:      0,
		lowSequence:  0,
	}
	return &butterfly
}

/*
	获取新的id
*/
func (b *Butterfly) Next() (int64, error) {
	b.Lock()
	if b.lowSequence == lowSequenceMax {
		if b.highSequence == highSequenceMax {
			if b.timeStamp == timestampMax {
				return 0, errors.New("no more id")
			} else {
				b.timeStamp++
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
	id := b.timeStamp<<timeStampShift | b.highSequence<<highSequenceShift | b.machine<<machineShift | b.lowSequence
	b.Unlock()
	return id, nil
}
