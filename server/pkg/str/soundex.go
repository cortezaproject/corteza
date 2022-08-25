package str

import (
	"strings"
)

// ToSoundex takes a word and returns the soundex code for it.
// https://en.wikipedia.org/wiki/Soundex
//
// 1. Retain the first letter of the name and drop all other occurrences of a, e, i, o, u, y, h, w.
// 2. Replace consonants with digits as follows (after the first letter):
//		b, f, p, v → 1
//		c, g, j, k, q, s, x, z → 2
//		d, t → 3
//		l → 4
//		m, n → 5
//		r → 6
// 3. If two or more letters with the same number are adjacent in the original name (before step 1),
//		only retain the first letter; also two letters with the same number separated
//		by 'h' or 'w' are coded as a single number, whereas such letters separated by a vowel are coded twice.
//		This rule also applies to the first letter.
// 4. Iterate the previous step until you have one letter and three numbers.
//		If you have too few letters in your word that you can't assign three numbers, append with zeros
//		until there are three numbers. If you have more than 3 letters, just retain the first 3 numbers.
func ToSoundex(s string) string {
	var (
		// soundex code
		code string
		// last code
		lastCode string
		// last rune
		lastRune rune
		// last rune is vowel
		lastRuneIsVowel bool
	)

	// retain the first letter of the name and drop all other occurrences of a, e, i, o, u, y, h, w
	for _, r := range s {
		if r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u' || r == 'y' || r == 'h' || r == 'w' {
			continue
		}

		code = string(r)
		break
	}

	// replace consonants with digits as follows (after the first letter)
	for _, r := range s {
		if r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u' || r == 'y' || r == 'h' || r == 'w' {
			lastRuneIsVowel = true
			continue
		}

		if lastRuneIsVowel {
			lastRuneIsVowel = false
			lastCode = ""
		}

		switch r {
		case 'b', 'f', 'p', 'v':
			lastCode = "1"
		case 'c', 'g', 'j', 'k', 'q', 's', 'x', 'z':
			lastCode = "2"
		case 'd', 't':
			lastCode = "3"
		case 'l':
			lastCode = "4"
		case 'm', 'n':
			lastCode = "5"
		case 'r':
			lastCode = "6"
		}

		if lastCode != "" && lastCode != string(lastRune) {
			code += lastCode
		}

		lastRune = r
	}

	// if two or more letters with the same number are adjacent in the original name (before step 1),
	// only retain the first letter
	// also two letters with the same number separated by 'h' or 'w' are coded as a single number,
	// whereas such letters separated by a vowel are coded twice
	// this rule also applies to the first letter
	code = strings.ReplaceAll(code, "11", "1")
	code = strings.ReplaceAll(code, "22", "2")
	code = strings.ReplaceAll(code, "33", "3")
	code = strings.ReplaceAll(code, "44", "4")
	code = strings.ReplaceAll(code, "55", "5")
	code = strings.ReplaceAll(code, "66", "6")

	// iterate the previous step until you have one letter and three numbers
	// if you have too few letters in your word that you can't assign three numbers,
	// append with zeros until there are three numbers
	// if you have more than 3 letters, just retain the first 3 numbers
	if len(code) < 4 {
		code += strings.Repeat("0", 4-len(code))
	} else {
		code = code[:4]
	}

	return code
}
