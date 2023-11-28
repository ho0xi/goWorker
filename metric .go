package worker

import "sync/atomic"

type Metric struct {
	busyWorkers uint64
}

func NewMetric() *Metric {
	return &Metric{
		busyWorkers: 0,
	}
}

// 增加计数
func (m *Metric) IncBusyWorker() uint64 {
	return atomic.AddUint64(&m.busyWorkers, 1)
}

// 减小计数
func (m *Metric) DecBusyWorker() uint64 {
	return atomic.AddUint64(&m.busyWorkers, ^uint64(0))
}

// 当前计数
func (m *Metric) BusyWorkers() uint64 {
	return atomic.LoadUint64(&m.busyWorkers)
}
