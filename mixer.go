package mixer

// Scalar implementation of Mixer

// **The Contract**
// - Mix to the shortest of dst, a and b
// - Does not allocate
// - Accepts any length
// - Returns values between -1 and 1

func MixScalar(dst, a, b []float32, gainA, gainB float32) {
	samples := findSampleCount(dst, a, b)
	for i := range samples {
		// Clamp(A * gainA + B * gainB) from -1 to 1
		dst[i] = clamp(a[i] * gainA + b[i] * gainB)
	}
}

func findSampleCount(dst, a, b []float32) int {
	n := len(dst)
	if len(a) < n {
		n = len(a)
	} 
	if len(b) < n {
		n = len(b)
	}
	return n
}

func clamp(mix float32) float32 {
	if mix > 1 {
		return 1
	}
	if mix < -1 {
		return -1
	}
	return mix
}