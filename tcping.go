package main

import (
	"fmt"
	"os"
	"tcping/utils"

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
		switch len(args) {
		case 0:
			cmd.Usage()
			return
		case 1:
			hostname := args[0]
			port := 80
			fmt.Println(hostname, port)
		case 2:
			hostname := args[0]
			port := args[1]
			fmt.Println(hostname, port)
		default:
			fmt.Println("invalid arguments!")
			fmt.Println()
			cmd.Usage()
		}

	},
}

func init() {
	// cmd.Flags().IntVarP(&port, "port", "p", 80, "port")
	rCmd.Flags().IntVarP(&counters, "counters", "c", utils.DCounters, "ping counter")
	rCmd.Flags().StringVarP(&interval, "interval", "i", "1s", "ping interval")
	rCmd.Flags().StringVarP(&timeout, "timeout", "t", "3s", "ping timeout")
}

func main() {
	if err := rCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
