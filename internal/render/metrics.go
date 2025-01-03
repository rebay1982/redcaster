package render

import (
	"time"
)

type fpsMetrics struct {
	startTime    time.Time
	metricsSum   int
	metricsIndex int
	metrics      [100]int
}

func (m *fpsMetrics) start() {
	m.startTime = time.Now()
}

func (m *fpsMetrics) stop() {
	delta := int(time.Since(m.startTime).Microseconds())

	m.update(delta)
}

func (m *fpsMetrics) update(delta int) {
	m.metricsSum += delta
	m.metricsSum -= m.metrics[m.metricsIndex]

	m.metrics[m.metricsIndex] = delta

	m.metricsIndex = (m.metricsIndex + 1) % len(m.metrics)
}

func (m fpsMetrics) getValue() float64 {
	avg := float64(m.metricsSum) / float64(len(m.metrics))
	return 1.0 / (0.000001 * avg)
}
