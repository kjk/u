package u

// StringArrayRemoveFirst removes first
func StringArrayRemoveFirst(a []string) []string {
	n := len(a)
	if n > 0 {
		copy(a[:n-1], a[1:])
		a = a[:n-1]
	}
	return a
}
