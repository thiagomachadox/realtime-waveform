package waveform

import (
	"math"
	"strings"
)

func GenerateASCII(samples []float32, width, height int) string {
	if len(samples) == 0 || width <= 0 || height <= 0 {
		return ""
	}

	// Compute RMS for each block
	blockSize := len(samples) / width
	rms := make([]float64, width)
	for i := 0; i < width; i++ {
		start := i * blockSize
		end := (i + 1) * blockSize
		if end > len(samples) {
			end = len(samples)
		}
		sum := 0.0
		for _, sample := range samples[start:end] {
			sum += float64(sample * sample)
		}
		rms[i] = math.Sqrt(sum / float64(blockSize))
	}

	// Normalize RMS values
	maxRMS := 0.0
	for _, v := range rms {
		if v > maxRMS {
			maxRMS = v
		}
	}
	for i := range rms {
		rms[i] /= maxRMS
	}

	// Generate ASCII art
	result := make([]string, height)
	for i := 0; i < height; i++ {
		row := make([]byte, width)
		for j, v := range rms {
			if int(v*float64(height)) >= height-i {
				row[j] = '#'
			} else {
				row[j] = ' '
			}
		}
		result[i] = string(row)
	}

	return strings.Join(result, "\n")
}