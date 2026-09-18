package mixer_test

import (
	"math"
	"testing"

	mixer "github.com/krish-bedi/simd-audio-kernel"
)

// Float32 can accurately represent 6 significant digits
// tolerance: 1e-6 for samples in range of [-1 to 1]

// Returns true if a-b is within tolerance
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
// go test -fuzz=FuzzMixScalar -fuzztime=30s
func FuzzMixScalar(f *testing.F) {
	// seed input
	f.Add(float32(0.25), float32(-0.5), float32(0.8), float32(0.6))

	f.Fuzz(func(t *testing.T, a, b, gainA, gainB float32) {
		// NaN is a valid float32 value
		if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) || 
			math.IsNaN(float64(gainA)) || math.IsNaN(float64(gainB)) {
			t.Skip()
		}

		dst := []float32{0}
		mixer.MixScalar(dst, []float32{a}, []float32{b}, gainA, gainB)
		if dst[0] < -1 || dst[0] > 1 {
			t.Fatalf("out of range: %v", dst[0])
		}
	})
}
