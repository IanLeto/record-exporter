package main

import (
	collector2 "exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	//collectors "record/collector"
)

func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Run(address string) {

}

var RootCmd = &cobra.Command{
	Use:   "tool", // 这个是命令的名字,跟使用没啥关系
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			collector prometheus.Collector
			err       error
			//dataSource collector2.DataSource
		)

		address, err := cmd.Flags().GetString("address")
		NoErr(err)
		kind, err := cmd.Flags().GetString("kind")
		NoErr(err)

		//c := cron.New()

		//_, err = c.AddFunc("*/65 * * * *", func() {
		//	err = dataSource.GetData()
		//	NoErr(err)
		//})
		NoErr(err)
		switch kind {
		case "ianRecord":
			collector = collector2.NewIanRecordCollector(address)
		case "filebeat":
			collector = collector2.NewFilebeatExporter(address)
		}
		//prometheus.MustRegister(collector)
		registry := prometheus.NewRegistry()
		registry.MustRegister(collector)
		http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		log.Fatal(http.ListenAndServe(":9101", nil))
	},
}

func init() {
	RootCmd.Flags().StringP("kind", "c", "", "config")
	RootCmd.Flags().StringP("address", "", "", "goOri ianRecord 访问方式")

}

func main() {
	NoErr(RootCmd.Execute())
}
