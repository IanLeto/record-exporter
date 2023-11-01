package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httputil"
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
	outputGaugeVec prometheus.GaugeVec
	client         http.Client
	address        string
}

func (f *FilebeatCollector) Describe(descs chan<- *prometheus.Desc) {
	f.outputGaugeVec.Describe(descs)
	f.beatGaugeVec.Describe(descs)
}

func (f *FilebeatCollector) Collect(metrics chan<- prometheus.Metric) {
	var (
		err error
		t   = &FilebeatResponse{}
	)
	request, err := http.NewRequest("GET", f.address, nil)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	resp, err := f.client.Do(request)
	if resp == nil {
		fmt.Println("Request failed:", err)
		f.outputGaugeVec.Collect(metrics)
		return
	}
	res, err := httputil.DumpResponse(resp, true)
	fmt.Printf("%s\n", res)
	err = json.Unmarshal(res, t)
	if err != nil {
		fmt.Println("Unmarshal failed:", err)
	}
	f.outputGaugeVec.WithLabelValues("events_acked").Set(float64(t.Libbeat.Output.Events.Acked))
	f.outputGaugeVec.WithLabelValues("events_active").Set(float64(t.Libbeat.Output.Events.Active))
	f.outputGaugeVec.WithLabelValues("events_batches").Set(float64(t.Libbeat.Output.Events.Batches))
	f.outputGaugeVec.WithLabelValues("events_total").Set(float64(t.Libbeat.Output.Events.Total))
	f.outputGaugeVec.Collect(metrics)
}

func NewFilebeatExporter(address string) *FilebeatCollector {
	return &FilebeatCollector{
		beatGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "filebeat",
			Help: "filebeat metrics",
		}, []string{"beat"}),
		outputGaugeVec: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "filebeat_output",
			Help: "filebeat output metrics",
		}, []string{"output"}),
	}
}
