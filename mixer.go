package mixer

// Scalar implementation of Mixer

// **The Contract**
// - Mix to the shortest of dst, a and b
// - Does not allocate
// - Accepts any length
// - Returns values between -1 and 1

func MixScalar(dst, a, b []float32, gainA, gainB float32) {
	sampleCount := findSmallestSlice(dst, a, b)
	for i := range sampleCount {
		// Clamp(A * gainA + B * gainB) from -1 to 1
		dst[i] = clamp(a[i] * gainA + b[i] * gainB)
	}
}

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