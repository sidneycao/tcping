package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/SidneyCao/tcping/utils"

	"github.com/spf13/cobra"
)

var (
	//port     int
	counters int
	interval string
	timeout  string
)

var rCmd = cobra.Command{
	Use:   "tcping [hostname/ip] [port]",
	Short: "The tcping is similarly to 'ping', but over tcp connection",
	Long:  "The tcping is similarly to 'ping', but over tcp connection, And written with Golang",
	Example: `
  1. ping www.google.com with default port 80
    > tcping www.google.com
  2. ping www.google.com with a custom port
    > tcping www.google.com 443
  `,
	Run: func(cmd *cobra.Command, args []string) {
		var host, port string
		switch len(args) {
		case 0:
			cmd.Usage()
			return
		case 1:
			host = args[0]
			port = "80"
		case 2:
			host = args[0]
			port = args[1]
		default:
			fmt.Println("invalid arguments!")
			fmt.Println()
			cmd.Usage()
			return
		}
		t := utils.NewTarget(host, port, timeout, interval)
		pinger := utils.NewPing(*t, os.Stdin, counters)
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
		fmt.Println("--- tcping starting ---")
		go pinger.Ping()
		select {
		case <-s:
		case <-pinger.Done():
		}
		pinger.Summarize()
	},
}

func init() {
	//rCmd.Flags().IntVarP(&port, "port", "p", 80, "port")
	rCmd.Flags().IntVarP(&counters, "counters", "c", utils.DCounters, "ping counter(if <= 0 will ping continuously)")
	rCmd.Flags().StringVarP(&interval, "interval", "i", "1", "ping interval, the unit is seconds")
	rCmd.Flags().StringVarP(&timeout, "timeout", "t", "3", "ping timeout, the unit is seconds")
}

func main() {
	if err := rCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
