package generator

import "time"

func NewGeneratorWithNowTime() *Butterfly {
	return NewButterfly(uint64(time.Now().UnixMilli()))
}
