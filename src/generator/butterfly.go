package generator

import (
	"sync"
)

const (
	timestampBits     = 41
	maxTimestamp      = -1 ^ (-1 << timestampBits)
	nodeBits          = 13                                            // 节点ID位数
	maxNodeID         = -1 ^ (-1 << nodeBits)                         // 节点ID最大值
	highSequenceBits  = 8                                             // 高序列号位数
	lowSequenceBits   = 1                                             // 低序列号位数
	maxHighSequence   = -1 ^ (-1 << highSequenceBits)                 // 高序列号最大值
	maxLowSequence    = -1 ^ (-1 << lowSequenceBits)                  // 低序列号最大值
	timeShift         = highSequenceBits + nodeBits + lowSequenceBits // 时间戳位移量
	highSequenceShift = nodeBits + lowSequenceBits                    // 低序列号位移量
	nodeIDShift       = lowSequenceBits                               // 节点ID位移量
)

type Butterfly struct {
	timestamp    uint64     // 时间戳
	highSequence uint64     // 高序列号
	lowSequence  uint64     // 低序列号
	nodeID       uint64     // 节点ID
	mutex        sync.Mutex // 互斥锁
}

func (b *Butterfly) Generate() uint64 {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.lowSequence = (b.lowSequence + 1) & maxLowSequence
	if b.lowSequence == 0 {
		b.nodeID = (b.nodeID + 1) & maxNodeID
		if b.nodeID == 0 {
			b.highSequence = (b.highSequence + 1) & maxHighSequence
			if b.highSequence == 0 {
				b.timestamp++
			}
		}
	}

	id := ((b.timestamp & maxTimestamp) << timeShift) |
		((b.highSequence & maxHighSequence) << highSequenceShift) |
		((b.nodeID & maxNodeID) << nodeIDShift) |
		(b.lowSequence & maxLowSequence)
	return id
}

func NewButterfly(initTimestamp uint64) *Butterfly {
	return &Butterfly{timestamp: initTimestamp}
}
