package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	collectors "record/collector"
)

func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

var RootCmd = &cobra.Command{
	Use:   "tool", // 这个是命令的名字,跟使用没啥关系
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		var collector prometheus.Collector
		address, err := cmd.Flags().GetString("address")
		NoErr(err)
		kind, err := cmd.Flags().GetString("kind")
		NoErr(err)
		switch kind {
		case "ianRecord":
			collector = collectors.NewIanRecordCollector(address)
		case "filebeat":
			collector = collectors.NewFilebeatExporter(address)
		default:
			collector = collectors.NewIanRecordCollector(address)
		}
		prometheus.MustRegister(collector)
		registry := prometheus.NewRegistry()
		registry.MustRegister(collector)
		http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		log.Fatal(http.ListenAndServe(":9101", nil))
	},
}

func init() {
	RootCmd.Flags().StringP("kind", "c", "", "config")
	RootCmd.Flags().StringP("address", "", "", "goOri ianRecord 访问方式")
	//RootCmd.Flags().BoolP("pass", "p", false, "pass")
	//RootCmd.Flags().Bool("debug", false, "debug")
	//RootCmd.Flags().String("init", "", "init db 啥的，要现保证各个依赖项，安装部署成功")

}

func main() {
	NoErr(RootCmd.Execute())
}
