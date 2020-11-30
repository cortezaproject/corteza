package expr

func Example_min() {
	eval(`min(2, 1, 3)`, nil)

	// output
	// float64(1)
}

func Example_max() {
	eval(`max(2, 1, 3)`, nil)

	// output
	// float64(3)
}

func Example_round() {
	eval(`round(3.14,1)`, nil)

	// output
	// 3.1
}

func Example_floor() {
	eval(`floor(3.14)`, nil)

	// output
	// float64(3)
}

func Example_ceil() {
	eval(`ceil(3.14)`, nil)

	// output
	// float64(4)
}
