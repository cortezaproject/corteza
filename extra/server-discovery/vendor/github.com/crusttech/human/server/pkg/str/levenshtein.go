package str

// write Levenshtein Distance search algorithm for strings
// https://en.wikipedia.org/wiki/Levenshtein_distance
func ToLevenshteinDistance(a, b string) int {
	var (
		// length of a
		la = len(a)
		// length of b
		lb = len(b)
		// distance matrix
		d = make([][]int, la+1)
	)

	// initialize distance matrix
	for i := 0; i <= la; i++ {
		d[i] = make([]int, lb+1)
		d[i][0] = i
	}

	for j := 0; j <= lb; j++ {
		d[0][j] = j
	}

	// calculate distance matrix
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			if a[i-1] == b[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				// fix this min function
				d[i][j] = min(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+1)
			}
		}
	}

	return d[la][lb]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	}

	if b < c {
		return b
	}

	return c
}
