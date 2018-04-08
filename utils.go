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
		if len(t) == 0 {
			continue
		}
		u[i] = t
		i++
	}
	u = u[:i]
	return u
}
