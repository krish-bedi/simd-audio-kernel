package audio

import (
	"errors"
)

var (
	ErrFrameMisconfig = errors.New("audio: frameSize and step must be greater than 0")
	ErrFrameStepGTEFramesize = errors.New("audio: step must be less than frameSize")
)

// Splits input into slices of frameSize, moving forward by step in each iteration
func Frames(samples []float32, frameSize, step int) ([][]float32, error) {
	// input sample is empty
	if len(samples) == 0 {
		return [][]float32{}, nil
	}

	// frameSize greater than length of samples
	if frameSize > len(samples) {
		return [][]float32{}, nil
	}

	// frameSize and step must be > 0
	if frameSize <= 0 || step <= 0 {
		return [][]float32{}, ErrFrameMisconfig
	}

	// step must be smaller than frameSize. Ensure samples aren't dropped
	if step >= frameSize {
		return [][]float32{}, ErrFrameStepGTEFramesize
	}

	// Number of frames that will be generated: ((n - frameSize) / step) + 1
	// 
	// (n - framesize): gives us the largest value the last complete frame can start with
	// / step : gives us the number of steps it will take to reach the previous frame
	// + 1: length is 1 + last index. i.e slice {0,1,2} has length 3
	frameCount := ((len(samples) - frameSize) / step) + 1

	frames := make([][]float32, frameCount)
	
	frameNum := 0
	i := 0
	for ; i+frameSize <= len(samples); i += step {
		frames[frameNum] = samples[i : i+frameSize]
		frameNum++
	}

	return frames, nil
}