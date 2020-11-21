package expr

func Example_trim() {
	eval(`trim(" foo ")`, nil)

	// output:
	// foo
}
func Example_trimLeft() {
	eval(`trimLeft(" foo ", " ")`, nil)

	// output:
	// foo
}
func Example_trimRight() {
	eval(`trimRight(" foo ", " ")`, nil)

	// output:
	// foo
}
func Example_toLower() {
	eval(`toLower("FOO")`, nil)

	// output:
	// foo
}
func Example_toUpper() {
	eval(`toUpper("foo")`, nil)

	// output:
	// FOO
}
func Example_shortest() {
	eval(`shortest("foo", "foobar")`, nil)

	// output:
	// foo
}
func Example_longest() {
	eval(`longest("foo", "foobar")`, nil)

	// output:
	// foobar
}
