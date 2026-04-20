package id

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sony/sonyflake"
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

var (
	initialized bool
	idQueue     = make(chan uint64, 8000)
	// Keep this to 1 for now as we need to enforce that any later ID is larger
	// than any previous one.
	// If we use different deviceIDs for routines, this isn't guaranteed to hold.
	makerCount = 1
)

func Init(ctx context.Context) {
	if initialized {
		return
	}
	initialized = true

	wg := &sync.WaitGroup{}
	wg.Add(makerCount)

	for i := 0; i < makerCount; i++ {
		makeIDer(ctx, wg, idQueue, uint64(i+1))
	}

	wg.Wait()

	// Give it some time to warm up with some ids
	// Arbitrarily picked a good timeout based on benchmark runs
	time.Sleep(time.Second * 1)
}

func makeIDer(ctx context.Context, wg *sync.WaitGroup, qq chan uint64, thr uint64) {
	go func() {
		ss, err := initSonyflake(thr)
		if err != nil {
			panic(fmt.Errorf("sonyflake init failed: %w", err))
		}

		wg.Done()
		for {
			select {
			case <-ctx.Done():
				return

			default:
				id, err := ss.NextID()
				if err != nil {
					panic(err)
				}
				qq <- id
			}
		}
	}()
}

func initSonyflake(thr uint64) (out *sonyflake.Sonyflake, err error) {
	settings := sonyflake.Settings{
		StartTime: time.Unix(jan1st2017, 0),
	}

	// check if SONYFLAKE_MACHINE_ID is set and create
	// MachineID function that returns it
	var id uint64
	if machineID := os.Getenv(envKeyMachineID); machineID != "" {
		// check if machineID is a valid uint16
		if id, err = strconv.ParseUint(machineID, 10, 16); err != nil {
			// should crash right away
			err = fmt.Errorf(
				"could not use %s (%s), expecting uint16 (0..65535): %w",
				envKeyMachineID,
				machineID,
				err,
			)
			return
		} else {
			settings.MachineID = func() (uint16, error) {
				return uint16(id + thr), nil
			}
		}
	} else if thr > 0 {
		settings.MachineID = func() (uint16, error) {
			return uint16(thr), nil
		}
	}

	out = sonyflake.NewSonyflake(settings)
	_, err = out.NextID()

	return
}

func Next() uint64 {
	if !initialized {
		panic("ID generator not initialized: call pkg/id.Init")
	}
	return <-idQueue
}
