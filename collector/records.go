package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

type IanRecordCollector struct {
	mealGaugeVec prometheus.GaugeVec
}

func (i *IanRecordCollector) Describe(descs chan<- *prometheus.Desc) {
	i.mealGaugeVec.Describe(descs)
}

func (i *IanRecordCollector) Collect(metrics chan<- prometheus.Metric) {
	i.mealGaugeVec.WithLabelValues().Set(1)
	i.mealGaugeVec.Collect(metrics)

}

func NewIanRecordCollector() *IanRecordCollector {
	return &IanRecordCollector{
		mealGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "ian_test_meal",
			Help: "help ian control the energy input",
		}, []string{}),
	}
}
