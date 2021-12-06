package expr

import (
	"time"
)

var (

	// human genome project end
	hgp = time.Date(2003, 4, 14, 0, 0, 0, 0, time.UTC)

	// world's first antibiotic
	wfa = time.Date(1928, 9, 128, 0, 0, 0, 0, time.UTC)

	// groundhog day
	ghd = time.Date(1993, 2, 2, 6, 0, 0, 0, time.FixedZone("", -5*60*60))

	exampleTimeParams = map[string]interface{}{
		"hgp": hgp,
		"wfa": wfa,
		"ghd": ghd,
		"now": now,
	}
)

func Example_strftime() {
	eval(`strftime(ghd, "%Y-%m-%dT%H:%M:%S")`, exampleTimeParams)

	// output:
	// 1993-02-02T06:00:00
}

func Example_strftimeWithModTime() {
	eval(`strftime(modTime(ghd, "+30m"), "%Y-%m-%dT%H:%M:%S")`, exampleTimeParams)

	// output:
	// 1993-02-02T06:30:00
}

func Example_parseISODate() {
	eval(`date("1993-02-02T06:00:00-05:00")`, nil)

	// output:
	// 1993-02-02 06:00:00 -0500 -0500
}

func Example_parseDate() {
	eval(`date("1993-02-02 06:00:00+01:10")`, nil)

	// output:
	// 1993-02-02 06:00:00 +0110 +0110
}

func Example_parseDuration() {
	eval(`parseDuration("2h")`, nil)

	// output:
	// 2h0m0s
}

func Example_earliest() {
	eval(`earliest(hgp, wfa)`, exampleTimeParams)

	// output:
	// 1929-01-06 00:00:00 +0000 UTC
}

func Example_latest() {
	eval(`latest(ghd, wfa)`, exampleTimeParams)

	// output:
	// 1993-02-02 06:00:00 -0500 -0500
}

func Example_isLeapYear() {
	eval(`isLeapYear(ghd)`, exampleTimeParams)

	// output:
	// false
}

func Example_isWeekDay() {
	eval(`isWeekDay(ghd)`, exampleTimeParams)

	// output:
	// true
}

func Example_modTime() {
	eval(`modTime(ghd, "+30m")`, exampleTimeParams)

	// output:
	// 1993-02-02 06:30:00 -0500 -0500
}

func Example_modWeek() {
	eval(`modWeek(ghd, "1")`, exampleTimeParams)

	// output:
	// 1993-02-09 06:00:00 -0500 -0500
}

func Example_modDate() {
	eval(`modDate(ghd, "5")`, exampleTimeParams)

	// output:
	// 1993-02-07 06:00:00 -0500 -0500
}

func Example_modMonth() {
	eval(`modMonth(ghd, "1")`, exampleTimeParams)

	// output:
	// 1993-03-02 06:00:00 -0500 -0500
}

func Example_modYear() {
	eval(`modYear(ghd, "5")`, exampleTimeParams)

	// output:
	// 1998-02-02 06:00:00 -0500 -0500
}

func Example_sub() {
	eval(`sub(hgp, ghd)`, exampleTimeParams)

	// output:
	// 321627600000
}

func Example_subErrorCheck() {
	eval(`sub(ghd, hgp)`, exampleTimeParams)

	// output:
	// error: can not evaluate sub(ghd, hgp): expecting 2nd input to be less than 1st input
}

func Example_unix() {
	eval(`strftime(ghd, "%L")`, exampleTimeParams)

	// output:
	// 728650800
}
