package router

// Holds a bunch of helper functions for dealing with labels.

// SplitDomainName splits a name string into it's labels.
// www.miek.nl. returns []string{"www", "miek", "nl"}
// The root label (.) returns nil. Note that using
// strings.Split(s) will work in most cases, but does not handle
// escaped dots (\.) for instance.
func splitDomainName(s string) (labels []string) {
	if len(s) == 0 {
		return nil
	}
	fqdnEnd := 0 // offset of the final '.' or the length of the name
	idx := split(s)
	begin := 0
	if s[len(s)-1] == '.' {
		fqdnEnd = len(s) - 1
	} else {
		fqdnEnd = len(s)
	}

	switch len(idx) {
	case 0:
		return nil
	case 1:
		// no-op
	default:
		end := 0
		for i := 1; i < len(idx); i++ {
			end = idx[i]
			labels = append(labels, s[begin:end-1])
			begin = end
		}
	}

	labels = append(labels, s[begin:fqdnEnd])
	return labels
}

// Split splits a name s into its label indexes.
// www.miek.nl. returns []int{0, 4, 9}, www.miek.nl also returns []int{0, 4, 9}.
// The root name (.) returns nil. Also see SplitDomainName.
func split(s string) []int {
	if s == "." {
		return nil
	}
	idx := make([]int, 1, 3)
	off := 0
	end := false

	for {
		off, end = nextLabel(s, off)
		if end {
			return idx
		}
		idx = append(idx, off)
	}
}

// nextLabel returns the index of the start of the next label in the
// string s starting at offset.
// The bool end is true when the end of the string has been reached.
// Also see PrevLabel.
func nextLabel(s string, offset int) (i int, end bool) {
	quote := false
	for i = offset; i < len(s)-1; i++ {
		switch s[i] {
		case '\\':
			quote = !quote
		default:
			quote = false
		case '.':
			if quote {
				quote = !quote
				continue
			}
			return i + 1, false
		}
	}
	return i + 1, true
}

// prevLabel returns the index of the label when starting from the right and
// jumping n labels to the left.
// The bool start is true when the start of the string has been overshot.
// Also see NextLabel.
func prevLabel(s string, n int) (i int, start bool) {
	if n == 0 {
		return len(s), false
	}
	lab := split(s)
	if lab == nil {
		return 0, true
	}
	if n > len(lab) {
		return 0, true
	}
	return lab[len(lab)-n], false
}
