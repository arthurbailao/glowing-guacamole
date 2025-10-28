package detector

import (
	"container/ring"
	"math"
)

// Detector is a detector of anomalies
type Detector struct {
	window            *ring.Ring
	maxWindowSize     int
	currentWindowSize int
	threshold         float64
}

// NewDetector creates a new detector
func NewDetector(threshold float64, windowSize int) *Detector {
	return &Detector{
		window:            ring.New(windowSize),
		maxWindowSize:     windowSize,
		currentWindowSize: 0,
		threshold:         threshold,
	}
}

// Detect returns the z-score and a boolean indicating if the value is an anomaly.
// It also returns a boolean indicating if it's warming up.
func (d *Detector) Detect(value float64) (float64, bool, bool) {
	d.add(value)

	if d.currentWindowSize < d.maxWindowSize {
		return 0, false, true
	}

	values := d.getRollingWindow()

	var sum float64
	for _, v := range values {
		sum += v
	}

	mean := sum / float64(len(values))

	var variance float64
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}

	if variance == 0 {
		return 0, false, false
	}

	stdDev := math.Sqrt(variance / float64(len(values)))
	z := math.Abs(value-mean) / stdDev

	return z, z >= d.threshold, false
}

func (d *Detector) add(value float64) {
	d.window.Value = value
	d.window = d.window.Next()
	d.currentWindowSize = min(d.currentWindowSize+1, d.maxWindowSize)
}

func (d *Detector) getRollingWindow() []float64 {
	var window = make([]float64, 0, d.currentWindowSize)

	d.window.Do(func(v any) {
		if v == nil {
			return
		}
		window = append(window, v.(float64))
	})

	return window
}
