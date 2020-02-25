package values

import (
	"regexp"
	"time"
)

const (
	strBoolTrue  = "1"
	strBoolFalse = "0"

	datetimeInputFormatDate = "2006-01-02"
	datetimeInputFormatTime = "15:04:05"
	datetimeInputFormatFull = time.RFC3339

	fieldOpt_Datetime_onlyDate         = "onlyDate"
	fieldOpt_Datetime_onlyTime         = "onlyTime"
	fieldOpt_Datetime_onlyFutureValues = "onlyFutureValues"
	fieldOpt_Datetime_onlyPastValues   = "onlyPastValues"

	fieldOpt_Number_precision = "precision"

	fieldOpt_Url_onlySecure = "onlySecure"
)

var (
	// value resembles something that can be true
	truthy = regexp.MustCompile(`^(t(rue)?|y(es)?|1)$`)

	// value resembles something that can be a reference
	refy = regexp.MustCompile(`^[1-9](\d*)$`)
)

func nowPtr() *time.Time {
	now := time.Now()
	return &now
}
