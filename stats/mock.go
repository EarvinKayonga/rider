package stats

import "time"

// Mock implementation for a statter.
type Mock struct{}

func (e Mock) Close() error {
	return nil
}

func (e Mock) Dec(stat string, value int64, rate float32) error {
	return nil
}
func (e Mock) Gauge(stat string, value int64, rate float32) error {
	return nil
}
func (e Mock) GaugeDelta(stat string, value int64, rate float32) error {
	return nil
}
func (e Mock) Inc(stat string, value int64, rate float32) error {
	return nil
}
func (e Mock) Raw(stat string, value string, rate float32) error {
	return nil
}
func (e Mock) Set(stat string, value string, rate float32) error {
	return nil
}
func (e Mock) SetInt(stat string, value int64, rate float32) error {
	return nil
}
func (e Mock) SetPrefix(prefix string) {}

func (e Mock) Timing(stat string, delta int64, rate float32) error {
	return nil
}
func (e Mock) TimingDuration(stat string, delta time.Duration, rate float32) error {
	return nil
}
