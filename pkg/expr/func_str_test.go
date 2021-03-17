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

func Example_title() {
	eval(`title("foo bAR")`, nil)

	// output:
	// Foo bAR
}

func Example_untitle() {
	eval(`untitle("Foo Bar")`, nil)

	// output:
	// foo Bar
}

func Example_repeat() {
	eval(`repeat("duran ", c)`, map[string]interface{}{"c": 2})

	// output:
	// duran duran
}

func Example_replace_all() {
	eval(`replace("foobar baz", "ba", "BA", c)`, map[string]interface{}{"c": -1})

	// output:
	// fooBAr BAz
}

func Example_replace_first() {
	eval(`replace("foobar baz", "ba", "BA", c)`, map[string]interface{}{"c": 1})

	// output:
	// fooBAr baz
}

func Example_isUrl() {
	eval(`isUrl("http:/example.tld")`, nil)

	// output:
	// false
}

func Example_isEmail() {
	eval(`isUrl("example.user+valid@example.tld")`, nil)

	// output:
	// true
}

func Example_split() {
	eval(`split("This will be split in:2-parts", ":")`, nil)

	// output:
	// [This will be split in 2-parts]
}

func Example_join() {
	eval(`join(exploded, ",")`, map[string][]string{"exploded": {"One", "two", "three"}})

	// output:
	// One,two,three
}

func Example_hasSubstring_caseS() {
	eval(`hasSubstring("foo BAR", "o b", true)`, nil)

	// output:
	// false
}

func Example_hasSubstring_caseI() {
	eval(`hasSubstring("foo BAR", "o b", false)`, nil)

	// output:
	// true
}

func Example_hasSuffix() {
	eval(`hasSuffix("foo bar", "ar")`, nil)

	// output:
	// true
}

func Example_hasPrefix() {
	eval(`hasPrefix("foo bar", "foo ")`, nil)

	// output:
	// true
}

func Example_shorten_word() {
	eval(`shorten("foo bar 1337, this is. test one - or three", "word", c)`, map[string]int{"c": 3})

	// output:
	// foo bar 1337 …
}

func Example_shorten_char() {
	eval(`shorten("foo bar 1337, this is. test one - or three", "char", c)`, map[string]int{"c": 22})

	// output:
	// foo bar 1337, this is …
}

func Example_camelize() {
	eval(`camelize("foo bar baz_test")`, nil)

	// output:
	// fooBarBaz_test
}

func Example_snakify() {
	eval(`snakify("foo bar baz_test")`, nil)

	// output:
	// foo_bar_baz_test
}

func Example_substring() {
	eval(`substring("foo bar baz_test", start, end)`, map[string]interface{}{"start": 2, "end": -1})

	// output:
	// o bar baz_test
}

func Example_substring_highStart() {
	eval(`substring("foo", start, end)`, map[string]interface{}{"start": 3, "end": -1})

	// output:
	//
}

func Example_substring_withEnd() {
	eval(`substring("foo bar baz_test", start, end)`, map[string]interface{}{"start": 2, "end": 4})

	// output:
	// o b
}

func Example_substring_endOverflow() {
	eval(`substring("foo bar baz_test", start, end)`, map[string]interface{}{"start": 2, "end": 100})

	// output:
	// o bar baz_test
}
