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