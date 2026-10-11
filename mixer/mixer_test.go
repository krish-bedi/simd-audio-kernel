package mixer_test

import (
	"fmt"
	"math"
	"testing"

	mixer "github.com/krish-bedi/simd-audio-kernel/mixer"
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

// SIMD Mixer implementation tests

func TestMixSIMDMatchesScalar(t *testing.T) {
	lanes := mixer.Lanes()
	lenghts := []int{0, 1, lanes - 1, lanes, lanes + 1, 2*lanes + 5}

	for _, n := range lenghts {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			a := make([]float32, n)
			b := make([]float32, n)
			mixScalar := make([]float32, n)
			mixSimd := make([]float32, n)

			for i := 0; i < n; i++ {
				a[i] = float32(i%21)/10 - 1
				b[i] = float32(i%51)/25 - 1
			}

			mixer.MixScalar(mixScalar, a, b, 0.6, 0.4)
			mixer.MixSIMD(mixSimd, a, b, 0.6, 0.4)

			for i := range mixScalar {
				if !withinTolerance(mixScalar[i], mixSimd[i]) {
					t.Fatalf("do not match: scalar=%v, simd=%v", mixScalar[i], mixSimd[i])
				}
			}
		})
	}
}

// Fuzz
// go test -fuzz=FuzzMixScalar -fuzztime=30s
// go test -fuzz=FuzzMixSIMD -fuzztime=30s

func FuzzMixScalar(f *testing.F) {
	fuzzMixImpl(f, mixer.MixScalar)
}

func FuzzMixSIMD(f *testing.F) {
	fuzzMixImpl(f, mixer.MixSIMD)
}

func fuzzMixImpl(f *testing.F, mix Mixer) {
	// seed input
	f.Add(float32(0.25), float32(-0.5), float32(0.8), float32(0.6))

	f.Fuzz(func(t *testing.T, a, b, gainA, gainB float32) {
		// NaN is a valid float32 value
		if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) || 
			math.IsNaN(float64(gainA)) || math.IsNaN(float64(gainB)) {
			t.Skip()
		}

		dst := []float32{0}
		mix(dst, []float32{a}, []float32{b}, gainA, gainB)
		if dst[0] < -1 || dst[0] > 1 {
			t.Fatalf("out of range: %v", dst[0])
		}
	})
}

// Benchmark

type Mixer func(dst, a, b []float32, gainA, gainB float32)

// Benchmark Scalar and SIMD Mixer implementations together and save result to file
// go test -run='^$' -bench='BenchmarkMix' -benchmem -count=10 ./... | tee bench.txt

func BenchmarkMix(b *testing.B) {
	sizes := map[string]int{
		"32": 32,
		"1K": 1<<10,
		"1M": 1<<20,
	}
	impls := map[string]Mixer{
		"scalar": mixer.MixScalar,
		"simd": mixer.MixSIMD,
	}

	for size, value := range sizes {
		for name, impl := range impls {
			b.Run(fmt.Sprintf("size=%s/impl=%s", size, name), func(b *testing.B) {
				benchmarkMixImpl(b, value, impl)
			})
		}
	}
}

// prevents compiler optimization from affecting benchmark timing
var benchmarkSink float32

func benchmarkMixImpl(b *testing.B, n int, mix Mixer) {
	trackA := make([]float32, n)
	trackB := make([]float32, n)
	dst := make([]float32, n)

	// deterministic repeating patterns let us reproduce the benchmark
	// math/rand with a fixed seed would work as well
	for i := range trackA {
		// populate both tracks with values from -1 to 1
		trackA[i] = float32(i%101) / 50 - 1
		trackB[i] = float32(i%21) / 10 - 1
	}

	// report throughput: 2 reads (trackA, trackB) + 1 write (dst)
	// each float32 (4 bytes)
	b.SetBytes(int64(n * 3 * 4))

	// Only measure time taken by the mixer
	b.ResetTimer()

	// timed loop
	for i := 0; i < b.N; i++ {
		mix(dst, trackA, trackB, 0.7, 0.3)
	}
	// save value from dst so compiler does not try to skip it during optimization
	benchmarkSink = dst[n-1]
}
