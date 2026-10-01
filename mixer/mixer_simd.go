package mixer

import (
	"simd"
)

func MixSIMD(dst, a, b []float32, gainA, gainB float32) {
	sampleCount := findSmallestSlice(dst, a, b)
	
	// Insert both gain and clamp values into each lane as they are
	// .. required for each operation
	vaGain := simd.BroadcastFloat32s(gainA)
	vbGain := simd.BroadcastFloat32s(gainB)
	lower := simd.BroadcastFloat32s(-1)
	upper := simd.BroadcastFloat32s(1)

	// Len() returns the # of lanes (length of vector)
	// decided by Go based on target hardware
	lanes := vaGain.Len()

	// Loop over the input tracks in counts of lanes
	// For ex: sampleCount = 10; lanes = 4
		// iter1 -> [a0, a1, a2, a3], [b0, b1, b2, b3]
		// iter2 -> [a4, a5, a6, a7], [b4, b5, b6, b7]
	// Process remaining [a8, a9], [b8, b9] in cleanup loop

	i := 0
	for ; i+lanes <= sampleCount; i+=lanes {
		va := simd.LoadFloat32s(a[i:])
		vb := simd.LoadFloat32s(b[i:])

		mixed := va.Mul(vaGain).Add(vb.Mul(vbGain))
		// Max returns the larger of the two values, Min returns the smaller
		// Max returns -1 for any value smaller than -1
		// Min returns +1 for any value bigger than +1
		clipped := mixed.Max(lower).Min(upper)
		clipped.Store(dst[i:])
	}

	// Cleanup loop handles remaining samples. iters = sampleCount % lanes
	// Scalar implementation is acceptable. At most lanes-1 elements
	for ; i < sampleCount; i++ {
		dst[i] = clamp(a[i]*gainA + b[i]*gainB)
	}

}