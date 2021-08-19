package expr

func Example_min() {
	eval(`min(2, 1, 3)`, nil)

	// output:
	// 1
}

func Example_max() {
	eval(`max(2, 1, 3)`, nil)

	// output:
	// 3
}

func Example_round() {
	eval(`round(3.14,1)`, nil)

	// output:
	// 3.1
}

func Example_floor() {
	eval(`floor(3.14)`, nil)

	// output:
	// 3
}

func Example_ceil() {
	eval(`ceil(3.14)`, nil)

	// output:
	// 4
}

func Example_abs() {
	eval(`abs(-3.14)`, nil)

	// output:
	// 3.14
}

func Example_log() {
	eval(`log(100)`, nil)

	// output:
	// 2
}

func Example_ln() {
	eval(`pow(2, 3)`, nil)

	// output:
	// 8
}

func Example_sqrt() {
	eval(`sqrt(16)`, nil)

	// output:
	// 4
}

func Example_sum() {
	eval(`sum(1, 2, "foo", 3.4, 4, "3")`, 10.4)

	// output:
	// 13.4
}

func Example_average() {
	eval(`average(1, 2, 3, "foo", 10)`, 10.4)

	// output:
	// 4
}

func Example_randomWithSingleInput() {
	eval(`max(random(6), 7)`, nil)

	// output:
	// 7
}

func Example_randomWithTwoInput() {
	eval(`min(random(2, 6), 1)`, nil)

	// output:
	// 1
}
