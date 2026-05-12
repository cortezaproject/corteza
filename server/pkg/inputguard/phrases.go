package inputguard

// instructionOverrides combines override phrases from all supported languages.
var instructionOverrides = combine(
	phrasesEN,
	phrasesES,
	phrasesFR,
	phrasesDE,
	phrasesPT,
	phrasesRU,
	phrasesZH,
	phrasesAR,
	phrasesJA,
	phrasesHI,
)

func combine(slices ...[]string) []string {
	total := 0
	for _, s := range slices {
		total += len(s)
	}
	out := make([]string, 0, total)
	for _, s := range slices {
		out = append(out, s...)
	}
	return out
}
