package mixer_test

import (
	"math"
	"testing"

	mixer "github.com/krish-bedi/simd-audio-kernel"
)

// Float32 can accurately represent 7 digits
// tolerance: 1e-6 for samples in range of [-1 to 1]
func withinTolerance(a, b float32) bool {
	const tolerance = 1e-6
	return math.Abs(float64(a-b)) <= tolerance
}

func TestMixScalar(t *testing.T) {
	a := []float32{0.2, 0.8, -0.8, 1.0}
	b := []float32{0.4, 0.8, -0.8, -1.0}
	// kernel does not allocate
	mix := make([]float32, len(a))

	mixer.MixScalar(mix, a, b, 0.5, 0.75)
	// Clamp(A * gainA + B * gainB, -1, 1)
	want := []float32{0.4, 1.0, -1.0, -0.25}
	
	for i := range want {
		if !withinTolerance(mix[i], want[i]) {
			t.Fatalf("sample %d, got %v, want %v", i, mix[i], want[i])
		}
	}
}

func TestMixScalarUsesShortestSlice(t *testing.T) {
	mix := []float32{9, 9, 9}
	mixer.MixScalar(mix, []float32{0.2, 0.4}, []float32{0.2}, 1, 1)

	want := []float32{0.4, 9, 9}
	for i := range want {
		if !withinTolerance(mix[i], want[i]) {
			t.Fatalf("sample %d, got %v, want %v", i, mix[i], want[i])
		}
	}
}
