package audio_test

import (
	"fmt"
	"slices"
	"testing"

	audio "github.com/krish-bedi/simd-audio-kernel/audio"
)

func TestFrames(t *testing.T) {
	samples := []float32{0,1,2,3,4,5,6,7,8,9,10}
	cases := []struct{
		frameSize int
		step int
		want [][]float32
		err error
	} {
		{4, 1, [][]float32{{0,1,2,3}, {1,2,3,4}, {2,3,4,5}, {3,4,5,6},
			{4,5,6,7}, {5,6,7,8}, {6,7,8,9}, {7,8,9,10}}, nil},
		{3, 2, [][]float32{{0,1,2}, {2,3,4}, {4,5,6}, {6,7,8}, {8,9,10}}, nil},
		{5, 3, [][]float32{{0,1,2,3,4}, {3,4,5,6,7}, {6,7,8,9,10}}, nil},
		// edge case: last frame is incomplete and should be dropped
		{4, 3, [][]float32{{0,1,2,3}, {3,4,5,6}, {6,7,8,9}}, nil},
		// frameSize > len(samples): should return empty slice as incomplete frames are dropped
		{12, 1, [][]float32{}, nil},
		// frameSize and step < 0: should return misconfig error
		{-1, -1, [][]float32{}, audio.ErrFrameMisconfig},
		// step >= frameSize: should return stepGTEFrameSize error
		{3, 3, [][]float32{}, audio.ErrFrameStepGTEFramesize},
		{3, 4, [][]float32{}, audio.ErrFrameStepGTEFramesize},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("frameSize: %d, step: %d", c.frameSize, c.step), func(t *testing.T) {
			got, err := audio.Frames(samples, c.frameSize, c.step)
			if err != c.err {
				t.Fatalf("expected error: %v, got: %v", c.err, err)
			}
			// Checks outer length and passes inner []float32 to slices.Equal
			slicesEqual := slices.EqualFunc(got, c.want, slices.Equal)
			if !slicesEqual {
				t.Errorf("expected slice: %v, got: %v", c.want, got)
			}
		})
	}
}
