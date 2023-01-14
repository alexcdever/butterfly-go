package model

type Butterfly struct {
	Base
	Bid          int64
	Timestamp    int64
	HighSequence int64
	Machine      int64
	LowSequence  int64
}

func (Butterfly) TableName() string {
	return "butterfly"
}
