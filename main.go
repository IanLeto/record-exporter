package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"log"
	"net/http"
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
		address, err := cmd.Flags().GetString("address")
		NoErr(err)
		Run(address)
	},
}

func init() {
	RootCmd.Flags().StringP("config", "c", "", "config")
	RootCmd.Flags().BoolP("pass", "p", false, "pass")
	RootCmd.Flags().Bool("debug", false, "debug")
	RootCmd.Flags().String("init", "", "init db 啥的，要现保证各个依赖项，安装部署成功")

}

func main() {

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
