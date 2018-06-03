package factory

import (
	"github.com/sony/sonyflake"
	"time"
)

// SonyflakeFactory is a configuration struct
type SonyflakeFactory struct {
	*sonyflake.Sonyflake
}

// Sonyflake is the active ID generator instance
var Sonyflake *SonyflakeFactory

func init() {
	Sonyflake = &SonyflakeFactory{
		sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: time.Unix(1503550784, 0),
		}),
	}
}

// NextID returns uint64 ID, escalates possible error to a panic
func (s *SonyflakeFactory) NextID() uint64 {
	// sonyflake errors out only when the time overflows, that will
	// occur in approximately 174 years after the custom epoch.
	// If the 10ms keyspace is exhausted, NextID will sleep and return
	// an ID from the next interval. It can't fail, because the generator
	// is protected by a Mutex.
	id, err := s.Sonyflake.NextID()
	if err != nil {
		panic(err)
	}
	return id
}
