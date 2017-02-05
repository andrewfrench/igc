package main

func sortCorrelations(list *[]Correlation) *[]Correlation {
	listLength := len(*list)

	if listLength < 2 {
		return list
	}

	left := new([]Correlation)
	right := new([]Correlation)

	for _, c := range (*list)[listLength / 2:] {
		*left = append(*left, c)
	}

	for _, c := range (*list)[:listLength / 2] {
		*right = append(*right, c)
	}

	left = sortCorrelations(left)
	right = sortCorrelations(right)

	return merge(left, right)
}

func merge(left, right *[]Correlation) *[]Correlation {
	sorted := new([]Correlation)

	for len(*left) > 0 && len(*right) > 0 {
		l := (*left)[0]
		r := (*right)[0]

		if l.Ratio > r.Ratio {
			*sorted = append(*sorted, l)
			*left = (*left)[1:]
		} else {
			*sorted = append(*sorted, r)
			*right = (*right)[1:]
		}
	}

	for _, c := range *left {
		*sorted = append(*sorted, c)
	}

	for _, c := range *right {
		*sorted = append(*sorted, c)
	}

	return sorted
}
