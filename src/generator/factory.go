package generator

import "time"

func NewGeneratorWithNowTime() *Butterfly {
	return NewButterfly(time.Now().UnixMilli())
}

func NewButterfly(initTimestamp int64) *Butterfly {
	return &Butterfly{timestamp: initTimestamp}
}

func NewButterflyList(timestamp int64) *ButterflyList {
	var list = &ButterflyList{
		generator:     NewButterfly(timestamp),
		AtLeastCount:  200,
		IncreaseCount: 3000,
	}
	list.construct()
	instanceList = append(instanceList, list)
	return list
}

func NewButterflyListWithNowTime() *ButterflyList {
	return NewButterflyList(time.Now().UnixMilli())
}
