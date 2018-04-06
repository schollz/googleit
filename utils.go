package googleit

// ListToSet convers a list to a set (removing duplicates)
func ListToSet(s []string) []string {
	m := make(map[string]struct{})
	for _, t := range s {
		m[t] = struct{}{}
	}
	u := make([]string, len(m))
	i := 0
	for t := range m {
		u[i] = t
		i++
	}
	return u
}
