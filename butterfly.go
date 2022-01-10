package main

import (
	"sync"
)

const (
	timeStampSize     int = 41
	highSequenceSize  int = 8
	machineSize       int = 13
	lowSequenceSize   int = 1
	timestampMax          = int64(-1 ^ (-1 << timeStampSize))
	highSequenceMax       = int64(-1 ^ (-1 << highSequenceSize))
	machineMax            = int64(-1 ^ (-1 << machineSize))
	lowSequenceMax        = int64(9)
	machineShift          = lowSequenceSize
	highSequenceShift     = machineShift + lowSequenceSize
	timeStampShift        = highSequenceShift + lowSequenceSize
)

type Butterfly struct {
	sync.Mutex
	timeStamp    int64
	highSequence int64
	machine      int64
	lowSequence  int64
}

func GetGenerator(timeStamp int64) Butterfly {
	butterfly := Butterfly{
		timeStamp:    timeStamp,
		highSequence: 0,
		machine:      0,
		lowSequence:  0,
	}
	return butterfly
}
func (b *Butterfly) Next() int64 {
	b.Lock()
	if b.lowSequence == lowSequenceMax {
		if b.highSequence == highSequenceMax {
			b.timeStamp++
			b.highSequence = 0
		} else {
			b.highSequence++
		}
		b.lowSequence = 0
	} else {
		b.lowSequence++
	}

	id := int64(b.timeStamp<<timeStampShift | b.highSequence<<highSequenceShift | b.machine<<machineShift | b.lowSequence)
	b.Unlock()
	return id
}
