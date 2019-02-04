package googleit

// ListToSet convers a list to a set (removing duplicates)
// but preserving order
func ListToSet(s []string) (t []string) {
	m := make(map[string]struct{})
	t = make([]string, len(s))
	i := 0
	for _, v := range s {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		t[i] = v
		i++
	}
	if i == 0 {
		return []string{}
	}
	t = t[:i]
	return
}
