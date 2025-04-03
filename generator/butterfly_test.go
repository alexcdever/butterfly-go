package generator

import (
	"testing"
	"time"
)

func TestButterfly_Generate(t *testing.T) {
	initTimestamp := time.Now().UnixNano() / int64(time.Millisecond) // 获取当前时间戳，单位毫秒
	b := NewButterfly(initTimestamp)

	// 测试生成的ID是否递增
	lastID := b.Generate()
	for i := 0; i < 1000; i++ {
		currentID := b.Generate()
		if currentID <= lastID {
			t.Errorf("ID not incrementing: %d, %d", currentID, lastID)
		}
		lastID = currentID
	}

	// 测试生成的ID是否符合预期
	b = NewButterfly(initTimestamp)
	expectedNodeID := int64(0)
	expectedTimestamp := initTimestamp
	expectedLowSequence := int64(1)
	expectedHighSequence := int64(0)
	for i := 0; i < 1000; i++ {
		id := b.Generate()

		nodeID := (id >> nodeIDShift) & maxNodeID
		if nodeID != expectedNodeID {
			t.Errorf("Unexpected node ID: %d, expected %d on %d times loop", nodeID, expectedNodeID, i)
		}

		timestamp := (id >> timeShift) & maxTimestamp
		if timestamp != expectedTimestamp {
			t.Errorf("Unexpected timestamp: %d, expected %d on %d times loop", timestamp, expectedTimestamp, i)
		}

		lowSequence := id & maxLowSequence
		if lowSequence != expectedLowSequence {
			t.Errorf("Unexpected low sequence: %d, expected %d on %d times loop", lowSequence, expectedLowSequence, i)
		}

		highSequence := (id >> highSequenceShift) & maxHighSequence
		if highSequence != expectedHighSequence {
			t.Errorf("Unexpected high sequence: %d, expected %d on %d times loop", highSequence, expectedHighSequence, i)
		}

		expectedLowSequence = (expectedLowSequence + 1) & maxLowSequence
		if expectedLowSequence == 0 {
			expectedNodeID = (expectedNodeID + 1) & maxNodeID
			if expectedNodeID == 0 {
				expectedHighSequence = (expectedHighSequence + 1) & maxHighSequence
				if expectedHighSequence == 0 {
					expectedTimestamp++
				}
			}
		}
	}
}

func TestButterflyList_Consume(t *testing.T) {
	initTimestamp := time.Now().UnixNano() / int64(time.Millisecond) // 获取当前时间戳，单位毫秒
	b := NewButterflyList(initTimestamp)
	consumeCount := 1000
	balanceAfterFirstConsume := len(b.UnusedIDList) - consumeCount - 1
	balanceAfterSecondConsume := len(b.UnusedIDList) - 2*consumeCount - 1

	//time.Sleep(time.Second * 10)

	if len(b.UnusedIDList) != b.IncreaseCount {
		t.Errorf("the lengh of unused id list expects as %d, but is %d", b.IncreaseCount, len(b.UnusedIDList))
	}

	// 测试生成的ID是否递增
	lastID := b.Consume()
	for i := 0; i < consumeCount; i++ {
		currentID := b.Consume()
		if currentID <= lastID {
			t.Errorf("ID not incrementing: %d, %d", currentID, lastID)
		}
		lastID = currentID
	}
	if len(b.UnusedIDList) != balanceAfterFirstConsume {
		t.Errorf("the lengh of unused id list expects as %d, but is %d", balanceAfterFirstConsume, len(b.UnusedIDList))
	}

	// 测试生成的ID是否递增
	idList := b.ConsumeInBatches(consumeCount)
	if len(b.UnusedIDList) != balanceAfterSecondConsume {
		t.Errorf("the lengh of unused id list expects as %d, but is %d", balanceAfterSecondConsume, len(b.UnusedIDList))
	}
	lastID = idList[0]
	for i := 1; i < len(idList); i++ {
		currentID := idList[i]
		if currentID <= lastID {
			t.Errorf("ID not incrementing: %d, %d", currentID, lastID)
		}
		lastID = currentID
	}
	if len(b.UnusedIDList) != balanceAfterSecondConsume {
		t.Errorf("the lengh of unused id list expects as %d, but is %d", balanceAfterSecondConsume, len(b.UnusedIDList))
	}

}
