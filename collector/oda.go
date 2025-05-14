package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"time"
)

var TransCountVec = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "TransCount",
	Help: "summary the trans count",
}, []string{"trans_type", "svc_name", "cluster", "trans_type_code", "ret_code", "pod_name", "project_name"})
var SuccessCountVec = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "SuccessCount",
	Help: "summary the success count",
}, []string{"trans_type", "svc_name", "cluster", "trans_type_code", "ret_code"})
var RespCountVec = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "RespCount",
	Help: "summary the resp count",
}, []string{"trans_type", "svc_name", "cluster", "trans_type_code", "ret_code"})

func NewDemoData() {
	go func() {
		for {
			select {
			case <-time.NewTicker(10 * time.Second).C:
				transType := fmt.Sprintf("trans_type_%d", rand.Intn(5)+1)
				svcName := fmt.Sprintf("cppas_%d", rand.Intn(3)+1)
				cluster := fmt.Sprintf("az%d", rand.Intn(2)+1)
				transTypeCode := fmt.Sprintf("AAAA%d", rand.Intn(10)+1)
				retCode := fmt.Sprintf("ret_code_%d", rand.Intn(3))
				TransCountVec.WithLabelValues(transType, svcName, cluster, transTypeCode, retCode, "CPAAS_TEST", "cpaas").Inc()
				TransCountVec.WithLabelValues(transType, svcName, cluster, transTypeCode, retCode, "Nrcp", "Nrcp").Inc()
				TransCountVec.WithLabelValues(transType, svcName, cluster, transTypeCode, retCode).Inc()
				SuccessCountVec.WithLabelValues(transType, svcName, cluster, transTypeCode, retCode).Inc()
				RespCountVec.WithLabelValues(transType, svcName, cluster, transTypeCode, retCode).Inc()
			}
		}
	}()
}
