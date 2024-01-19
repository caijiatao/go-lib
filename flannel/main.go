package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flannelFlags := flag.NewFlagSet("flannel", flag.ExitOnError)

	var kubeSubMgr bool

	flannelFlags.BoolVar(&kubeSubMgr, "kube", false, "test")

	err := flannelFlags.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("kube:", kubeSubMgr)
}
