package id

import (
	"github.com/sony/sonyflake"
	"time"
)

// Midnight, January 1st 2017
const jan1st2017 = 1483228800

var _sonyflake *sonyflake.Sonyflake

func init() {
	_sonyflake = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Unix(jan1st2017, 0),
	})
}

// NextID returns uint64 ID, or panics
//
// See https://github.com/sony/sonyflake for details
func nextSonyflake() uint64 {
	if id, err := _sonyflake.NextID(); err != nil {
		panic(err)
	} else {
		return id
	}
}

func Next() uint64 {
	return nextSonyflake()
}
