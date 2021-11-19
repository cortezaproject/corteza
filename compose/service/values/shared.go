package values

import (
	"regexp"
	"time"
)

const (
	strBoolTrue     = "1"
	strBoolFalse    = "0"
	strBoolFalseAlt = ""

	datetimeInternalFormatDate = "2006-01-02"
	datetimeIntenralFormatTime = "15:04:05"
	datetimeInternalFormatFull = time.RFC3339

	fieldOpt_Datetime_onlyDate         = "onlyDate"
	fieldOpt_Datetime_onlyTime         = "onlyTime"
	fieldOpt_Datetime_onlyFutureValues = "onlyFutureValues"
	fieldOpt_Datetime_onlyPastValues   = "onlyPastValues"

	fieldOpt_Url_onlySecure = "onlySecure"
)

var (
	// value resembles something that can be true
	truthy = regexp.MustCompile(`^(t(rue)?|y(es)?|1)$`)

	// value resembles something that can be a reference
	refy = regexp.MustCompile(`^[1-9](\d*)$`)

	// value resembels something that can be parsed as ISO 8601
	isoDaty     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`)
	hasTimezone = regexp.MustCompile(`(Z|\+\d{2}:\d{2})$`)
)
