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
	timestamp    int64      // 时间戳
	highSequence int64      // 高序列号
	lowSequence  int64      // 低序列号
	nodeID       int64      // 节点ID
	mutex        sync.Mutex // 互斥锁
}

func (b *Butterfly) Generate() int64 {
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

func (b *Butterfly) GenerateInBatches(count int) []int64 {
	var idList []int64
	for i := 0; i < count; i++ {
		idList = append(idList, b.Generate())
	}
	return idList
}

type ButterflyList struct {
	generator    *Butterfly
	mutex        sync.Mutex
	UnusedIDList []int64
	Period       int64
	// AtLeastCount required the length of UnusedIDList at least 200
	AtLeastCount int `validate:"required,gte=200"`
	// IncreaseCount required the count of appending to UnusedIDList at least 3000
	IncreaseCount int `validate:"required,gte=3000"`
}

func (b *ButterflyList) construct() {
	b.UnusedIDList = append(b.UnusedIDList, b.generator.GenerateInBatches(b.IncreaseCount)...)
}

func (b *ButterflyList) Consume() int64 {
	b.mutex.Lock()
	id := b.UnusedIDList[0]
	b.UnusedIDList = b.UnusedIDList[1:]
	b.mutex.Unlock()

	return id
}

func (b *ButterflyList) ConsumeInBatches(count int) (idList []int64) {
	b.mutex.Lock()
	idList = b.UnusedIDList[:count]
	b.UnusedIDList = b.UnusedIDList[count:]
	b.mutex.Unlock()
	return idList
}
