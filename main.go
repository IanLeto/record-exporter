package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	collector2 "record/collector"
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
		var (
			// 一组collectors，用来将获取需要的数据
			collectors []prometheus.Collector
			err        error
			// 意味着，我们对server 暴露出来的是http协议
			opts = promhttp.HandlerOpts{
				EnableOpenMetrics: false,
			}
			// 注册器，用来管理collectors 和 指标本身
			// 当请求打过来，例如http，会遍历collectors，然后调用collect方法，将数据写入到metrics中
			registry = prometheus.NewRegistry()
		)

		address, err := cmd.Flags().GetString("address")
		NoErr(err)
		kind, err := cmd.Flags().GetString("kind")
		NoErr(err)
		switch kind {
		case "ianRecord":
			collectors = append(collectors, collector2.NewIanRecordCollector(address))
		case "filebeat":
			//collector = collector2.NewFilebeatExporter(address)
		default:
			collectors = append(collectors, collector2.NewIanRecordCollector(address))
		}
		registry.MustRegister(collectors...)
		//registry.MustRegister(collector2.TransCountVec, collector2.SuccessCountVec, collector2.RespCountVec)
		// 注册一个handler，用来处理metrics请求
		// http 过来时候，会调用gather， reg 本身实现了gather，gather 会遍历所有的collectors 然后调用collect方法
		http.Handle("/metrics", promhttp.HandlerFor(registry, opts))
		collector2.NewDemoData()
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
