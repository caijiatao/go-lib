package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfgFile string
	port    int
)

func runKubelet(cmd *cobra.Command, args []string) error {
	fmt.Println(args)
	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "kubelet",
		Long: "kubelet is the primary node agent that runs on each node. this is a long desc",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runKubelet(cmd, args); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
	fmt.Println(cfgFile)
	fmt.Println(port)
	return cmd
}

func main() {
	cmd := NewCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
