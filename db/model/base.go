package model

import (
	"time"
)

type Base struct {
	CreatedDate time.Time
	CreatedTime time.Time
	UpdatedDate time.Time
	UpdatedTime time.Time
	DeletedDate time.Time
	DeletedTime time.Time
}
