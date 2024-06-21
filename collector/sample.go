package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	// 定义一个新的指标，这里是一个计数器
	SampleCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sample_metric_total",
			Help: "A sample counter metric",
		},
	)
)
