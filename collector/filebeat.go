package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"net/http"
)

type FilebeatResponse struct {
	Beat struct {
		Cgroup struct {
			Memory struct {
				Mem struct {
					Usage struct {
						Bytes int `json:"bytes"`
					} `json:"usage"`
				} `json:"mem"`
			} `json:"memory"`
		} `json:"cgroup"`
		Cpu struct {
			System struct {
				Ticks int `json:"ticks"`
				Time  struct {
					Ms int `json:"ms"`
				} `json:"time"`
			} `json:"system"`
			Total struct {
				Ticks int `json:"ticks"`
				Time  struct {
					Ms int `json:"ms"`
				} `json:"time"`
				Value int `json:"value"`
			} `json:"total"`
			User struct {
				Ticks int `json:"ticks"`
				Time  struct {
					Ms int `json:"ms"`
				} `json:"time"`
			} `json:"user"`
		} `json:"cpu"`
		Handles struct {
			Limit struct {
				Hard int `json:"hard"`
				Soft int `json:"soft"`
			} `json:"limit"`
			Open int `json:"open"`
		} `json:"handles"`
		Info struct {
			EphemeralId string `json:"ephemeral_id"`
			Uptime      struct {
				Ms int `json:"ms"`
			} `json:"uptime"`
			Version string `json:"version"`
		} `json:"info"`
		Memstats struct {
			GcNext      int   `json:"gc_next"`
			MemoryAlloc int   `json:"memory_alloc"`
			MemoryTotal int64 `json:"memory_total"`
			Rss         int   `json:"rss"`
		} `json:"memstats"`
		Runtime struct {
			Goroutines int `json:"goroutines"`
		} `json:"runtime"`
	} `json:"beat"`
	Filebeat struct {
		Events struct {
			Active int `json:"active"`
			Added  int `json:"added"`
			Done   int `json:"done"`
		} `json:"events"`
		Harvester struct {
			OpenFiles int `json:"open_files"`
			Running   int `json:"running"`
			Started   int `json:"started"`
		} `json:"harvester"`
	} `json:"filebeat"`
	Libbeat struct {
		Config struct {
			Module struct {
				Running int `json:"running"`
			} `json:"module"`
		} `json:"config"`
		Output struct {
			Events struct {
				Acked   int `json:"acked"`
				Active  int `json:"active"`
				Batches int `json:"batches"`
				Total   int `json:"total"`
			} `json:"events"`
			Read struct {
				Bytes int `json:"bytes"`
			} `json:"read"`
			Write struct {
				Bytes int `json:"bytes"`
			} `json:"write"`
		} `json:"output"`
		Pipeline struct {
			Clients int `json:"clients"`
			Events  struct {
				Active    int `json:"active"`
				Filtered  int `json:"filtered"`
				Published int `json:"published"`
				Total     int `json:"total"`
			} `json:"events"`
			Queue struct {
				Acked int `json:"acked"`
			} `json:"queue"`
		} `json:"pipeline"`
	} `json:"libbeat"`
	Registrar struct {
		States struct {
			Current int `json:"current"`
			Update  int `json:"update"`
		} `json:"states"`
		Writes struct {
			Success int `json:"success"`
			Total   int `json:"total"`
		} `json:"writes"`
	} `json:"registrar"`
	System struct {
		Load struct {
			Field1 float64 `json:"1"`
			Field2 float64 `json:"5"`
			Field3 float64 `json:"15"`
			Norm   struct {
				Field1 float64 `json:"1"`
				Field2 float64 `json:"5"`
				Field3 float64 `json:"15"`
			} `json:"norm"`
		} `json:"load"`
	} `json:"system"`
}

type FilebeatCollector struct {
	beatGaugeVec   prometheus.GaugeVec
	eventGaugeVec  prometheus.GaugeVec
	handleGaugeVec prometheus.GaugeVec
	harvesterGauge prometheus.GaugeVec

	outputGaugeVec prometheus.GaugeVec
	pipelineGauge  prometheus.GaugeVec

	up      *prometheus.Desc
	client  http.Client
	address string
	current float64
}

// Describe 描述所有我们想要收集的指标，尤其指标名，标签以及帮助信息
func (f *FilebeatCollector) Describe(descs chan<- *prometheus.Desc) {
	f.eventGaugeVec.Describe(descs)
	f.beatGaugeVec.Describe(descs)
	f.outputGaugeVec.Describe(descs)
	f.harvesterGauge.Describe(descs)
	f.pipelineGauge.Describe(descs)
	f.handleGaugeVec.Describe(descs)
}

func (f *FilebeatCollector) Collect(metrics chan<- prometheus.Metric) {
	var (
		err error
		t   = &FilebeatResponse{}
	)
	request, err := http.NewRequest("GET", f.address, nil)
	if err != nil {
		fmt.Println("Request failed:", err)
		metrics <- prometheus.MustNewConstMetric(f.up, prometheus.GaugeValue, 0) // 需要将up指标设置为0
		return
	}
	resp, err := f.client.Do(request)
	if resp == nil {
		fmt.Println("Request failed:", err)
		metrics <- prometheus.MustNewConstMetric(f.up, prometheus.GaugeValue, 0)
		//f.eventGaugeVec.Collect(metrics)
		return
	}
	//res, err := httputil.DumpResponse(resp.Body, true)
	res, err := io.ReadAll(resp.Body)

	fmt.Printf("%s\n", res)
	err = json.Unmarshal(res, t)
	if err != nil {
		fmt.Println("Unmarshal failed:", err)
	}

	f.eventGaugeVec.WithLabelValues("events_acked").Set(float64(t.Filebeat.Events.Active))
	f.eventGaugeVec.WithLabelValues("events_active").Set(float64(t.Filebeat.Events.Added))
	f.eventGaugeVec.WithLabelValues("events_batches").Set(float64(t.Filebeat.Events.Done))

	f.eventGaugeVec.WithLabelValues("events_read_bytes").Set(float64(t.Libbeat.Output.Read.Bytes))
	f.eventGaugeVec.WithLabelValues("events_write_bytes").Set(float64(t.Libbeat.Output.Write.Bytes))

	f.handleGaugeVec.WithLabelValues("handle_open_files").Set(float64(t.Beat.Handles.Open))
	f.handleGaugeVec.WithLabelValues("handle_limit_hard").Set(float64(t.Beat.Handles.Limit.Hard))
	f.handleGaugeVec.WithLabelValues("handle_limit_soft").Set(float64(t.Beat.Handles.Limit.Soft))

	f.harvesterGauge.WithLabelValues("harvester_open_files").Set(float64(t.Filebeat.Harvester.OpenFiles))
	f.harvesterGauge.WithLabelValues("harvester_running").Set(float64(t.Filebeat.Harvester.Running))
	f.harvesterGauge.WithLabelValues("harvester_started").Set(float64(t.Filebeat.Harvester.Started))

	f.outputGaugeVec.WithLabelValues("output_events_acked").Set(float64(t.Libbeat.Output.Events.Acked))
	f.outputGaugeVec.WithLabelValues("output_events_active").Set(float64(t.Libbeat.Output.Events.Active))
	f.outputGaugeVec.WithLabelValues("output_events_batches").Set(float64(t.Libbeat.Output.Events.Batches))
	f.outputGaugeVec.WithLabelValues("output_events_total").Set(float64(t.Libbeat.Output.Events.Total))

	f.pipelineGauge.WithLabelValues("pipeline_clients").Set(float64(t.Libbeat.Pipeline.Clients))
	f.pipelineGauge.WithLabelValues("pipeline_events_active").Set(float64(t.Libbeat.Pipeline.Events.Active))
	f.pipelineGauge.WithLabelValues("pipeline_events_filtered").Set(float64(t.Libbeat.Pipeline.Events.Filtered))
	f.pipelineGauge.WithLabelValues("pipeline_events_published").Set(float64(t.Libbeat.Pipeline.Events.Published))
	f.pipelineGauge.WithLabelValues("pipeline_events_total").Set(float64(t.Libbeat.Pipeline.Events.Total))
	f.pipelineGauge.WithLabelValues("pipeline_queue_acked").Set(float64(t.Libbeat.Pipeline.Queue.Acked))

	f.eventGaugeVec.Collect(metrics)

	f.beatGaugeVec.Collect(metrics)
	f.outputGaugeVec.Collect(metrics)
	f.harvesterGauge.Collect(metrics)
	f.pipelineGauge.Collect(metrics)
	f.handleGaugeVec.Collect(metrics)
}

func NewFilebeatExporter(address string) *FilebeatCollector {
	return &FilebeatCollector{
		address: address,
		beatGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "filebeat",
			Help: "filebeat metrics",
		}, []string{"beat"}),
		eventGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "filebeat_event",
			Help: "filebeat 输出的相关信息",
		}, []string{"event"}),
		handleGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "filebeat_handle",
			Help: "filebeat handle 相关信息",
		}, []string{"handle"}),
		harvesterGauge: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "filebeat_harvester",
			Help: "filebeat harvester 相关信息",
		}, []string{"harvester"}),

		pipelineGauge: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "filebeat_pipeline",
			Help: "filebeat pipeline 相关信息",
		}, []string{"pipeline"}),

		outputGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "filebeat_output",
			Help: "filebeat output 相关信息",
		}, []string{"output"}),

		up: prometheus.NewDesc("filebeat_up", "filebeat up", nil, nil),
	}
}
