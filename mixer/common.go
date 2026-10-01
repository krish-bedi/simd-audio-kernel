package mixer

func findSmallestSlice(dst, a, b []float32) int {
	n := len(dst)
	if len(a) < n {
		n = len(a)
	} 
	if len(b) < n {
		n = len(b)
	}
	return n
}

func clamp(n float32) float32 {
	if n > 1 {
		return 1
	}
	if n < -1 {
		return -1
	}
	return n
}