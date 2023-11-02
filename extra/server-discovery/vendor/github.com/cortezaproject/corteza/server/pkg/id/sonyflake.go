package id

import (
	"fmt"
	"github.com/sony/sonyflake"
	"os"
	"strconv"
	"time"
)

const (
	// Midnight, January 1st 2017
	jan1st2017 = 1483228800

	// Env key for custom machine ID
	//
	// Sonyflake assumes that your private IP falls under the following ranges:
	//     10.0.0.0 - 10.255.255.255
	//   172.16.0.0 - 172.31.255.255
	//  192.168.0.0 - 192.168.255.255
	//
	// When running in a container, this could not always the case.
	// To override this, set SONYFLAKE_MACHINE_ID to a valid uint16
	// (0..65535) and sonyflake will use that as machine ID.
	envKeyMachineID = "SONYFLAKE_MACHINE_ID"

	// When set to uint32 value, we no longer use sonyflake
	// but a simple atomic counter
	//
	// This is useful for testing and debugging
	// and should not be used in production
	// @todo implement
	// envKeySimpleIncrement = "SONYFLAKE_SIMPLE_INCREMENT"
)

var sf *sonyflake.Sonyflake

func init() {
	if err := initSonyflake(); err != nil {
		panic(fmt.Errorf("sonyflake init failed: %w", err))
	}
}

func initSonyflake() error {
	settings := sonyflake.Settings{
		StartTime: time.Unix(jan1st2017, 0),
	}

	// check if SONYFLAKE_MACHINE_ID is set and create
	// MachineID function that returns it
	if machineID := os.Getenv(envKeyMachineID); machineID != "" {
		// check if machineID is a valid uint16
		if id, err := strconv.ParseUint(machineID, 10, 16); err != nil {
			// should crash right away
			return fmt.Errorf(
				"could not use %s (%s), expecting uint16 (0..65535): %w",
				envKeyMachineID,
				machineID,
				err,
			)
		} else {
			settings.MachineID = func() (uint16, error) {
				return uint16(id), nil
			}
		}
	}

	sf = sonyflake.NewSonyflake(settings)
	_, err := sf.NextID()
	return err
}

// NextID returns uint64 ID, or panics
//
// See https://github.com/sony/sonyflake for details
func nextSonyflake() uint64 {
	if id, err := sf.NextID(); err != nil {
		panic(err)
	} else {
		return id
	}
}

func Next() uint64 {
	return nextSonyflake()
}
