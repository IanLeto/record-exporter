package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"net/http"
	"time"
)

type IanRecordCollector struct {
	mealGaugeVec   prometheus.GaugeVec
	costCountVec   prometheus.CounterVec
	weightGaugeVec prometheus.GaugeVec
	client         http.Client
	address        string
	interval       int
}

type TRecord struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	NorWeight  float64 `json:"mor_weight"`
	NigWeight  float64 `json:"nig_weight"`
	Cost       int     `json:"cost"`
	UpdateTime int64   `json:"create_time"` // 注意：JSON 字段是 create_time
}

type RecordResponse struct {
	Data struct {
		Items []TRecord `json:"items"`
	} `json:"Data"`
	Message   string `json:"Message"`
	TransType int    `json:"TransType"`
}

func (i *IanRecordCollector) GetData() error {
	return nil
}

func getTodayStartEnd() (int64, int64) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)
	return start.Unix(), end.Unix()
}

// Collect 收集指标
func (i *IanRecordCollector) Collect(metrics chan<- prometheus.Metric) {
	var resData = &RecordResponse{}

	defer func() {
		i.mealGaugeVec.Collect(metrics)
		i.weightGaugeVec.Collect(metrics)
		i.costCountVec.Collect(metrics)
	}()

	start, end := getTodayStartEnd()
	url := fmt.Sprintf("%s/v1/record?region=win&start_time=%d&end_time=%d", i.address, start, end)

	fmt.Println("📡 Sending request:")
	fmt.Println("  ➜ URL:", url)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("❌ Failed to create request:", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := i.client.Do(request)
	if err != nil {
		fmt.Println("❌ HTTP request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("❌ Failed to read response body:", err)
		return
	}

	fmt.Println("📥 Raw response body:")
	fmt.Println(string(body))

	err = json.Unmarshal(body, resData)
	if err != nil {
		fmt.Println("❌ Failed to unmarshal JSON:", err)
		return
	}

	if len(resData.Data.Items) == 0 {
		fmt.Println("⚠️ No items found in response.")
		return
	}

	// 只取第一个 item
	t := resData.Data.Items[0]
	fmt.Printf("✅ Parsed First Record: %+v\n", t)

	//period := TimePeriod(t.UpdateTime)
	i.weightGaugeVec.WithLabelValues("nor").Set(t.NorWeight)
	i.weightGaugeVec.WithLabelValues("nig").Set(t.NigWeight)
	i.costCountVec.WithLabelValues("nor").Add(float64(t.Cost))
	i.costCountVec.WithLabelValues("nig").Add(float64(t.Cost))
}

// Describe 向prometheus注册指标，他描述了我们想要收集的指标的名字，标签和帮助信息；
// 当collector 被注册到registry中时，会调用这个方法
func (i *IanRecordCollector) Describe(descs chan<- *prometheus.Desc) {
	i.mealGaugeVec.Describe(descs)
	i.costCountVec.Describe(descs)
	i.weightGaugeVec.Describe(descs)
}

func NewIanRecordCollector(address string) *IanRecordCollector {
	return &IanRecordCollector{
		mealGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:        "ian_test_meal",
			Help:        "help ian control the energy input",
			ConstLabels: prometheus.Labels{},
		}, []string{"time_period"}), // 定义一个时间段标签

		costCountVec: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "sum_money",
			Help: "Ian's sum money records",
		}, []string{"time_period"}), // 定义一个时间段标签

		weightGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weight",
			Help: "Ian's weight records",
		}, []string{"time_period"}), // 定义一个时间段标签

		client:  http.Client{},
		address: address,
	}
}

type TRecordToMetrics struct {
}

var Weight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "weight",
	Help: "Ian's weight records, divided into morning, afternoon, and evening measurements",
}, []string{"time"})

var SumMoney = *prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "sum_money",
	Help: "Ian's sum money records",
}, []string{"time"})
