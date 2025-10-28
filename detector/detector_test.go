package detector_test

import (
	"testing"

	"github.com/arthurbailao/cint-ad/detector"
	"github.com/stretchr/testify/assert"
)

func TestDetectZCalculation(t *testing.T) {
	d := detector.NewDetector(1.9, 5)
	for _, v := range []float64{10, 11, 12, 13} {
		d.Detect(v)
	}

	z, _, _ := d.Detect(10)
	assert.InDelta(t, 1.029, z, 0.001)
}

func TestDetectNoAnomalies(t *testing.T) {
	testCases := []struct {
		name   string
		values []float64
	}{
		{name: "less values than window size", values: []float64{10, 11, 12, 9999}},
		{name: "exceeds window size", values: []float64{10, 11, 12, 13, 10, 9}},
		{name: "negative numbers", values: []float64{-2, -1, 0, 1, 2, 3}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := detector.NewDetector(3, 5)
			for _, v := range tc.values {
				_, anomaly, _ := d.Detect(v)
				assert.False(t, anomaly)
			}
		})
	}
}

func TestDetectAnomalies(t *testing.T) {
	testCases := []struct {
		name     string
		values   []float64
		expected []bool
	}{
		{
			name:     "anomaly detected for new value",
			values:   []float64{10, 11, 12, 13, 9999},
			expected: []bool{false, false, false, false, true},
		},
		{
			name:     "anomaly detected for new negative value",
			values:   []float64{10, 11, 12, 13, -10},
			expected: []bool{false, false, false, false, true},
		},
		{
			name:     "anomaly detected",
			values:   []float64{10, 11, 12, 13, 12, 100, 12, 900},
			expected: []bool{false, false, false, false, false, true, false, true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := detector.NewDetector(1.9, 5)

			var result []bool
			for _, v := range tc.values {
				_, anomaly, _ := d.Detect(v)
				result = append(result, anomaly)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDetectNoAnamoliesEqualValues(t *testing.T) {
	d := detector.NewDetector(1.9, 5)
	values := []float64{10, 10, 10, 10, 10, 10}
	for _, v := range values {
		z, anomaly, _ := d.Detect(v)
		assert.Zero(t, z)
		assert.False(t, anomaly)
	}
}

func TestDetectWarmingUp(t *testing.T) {
	d := detector.NewDetector(1.9, 5)
	values := []float64{10, 9, 8, 7}
	for _, v := range values {
		z, anomaly, warming := d.Detect(v)
		assert.Zero(t, z)
		assert.False(t, anomaly)
		assert.True(t, warming)
	}

	_, anomaly, warming := d.Detect(8)
	assert.False(t, anomaly)
	assert.False(t, warming)
}
