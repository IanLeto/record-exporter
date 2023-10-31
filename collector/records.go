package collector

import "github.com/prometheus/client_golang/prometheus"

type IanRecordCollector struct {
	mealGaugeVec prometheus.GaugeVec
}

func NewIanRecordCollector() *IanRecordCollector {
	return &IanRecordCollector{
		mealGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "ian_test_meal",
			Help: "hei",
		}, []string{"meal"}),
	}
}
