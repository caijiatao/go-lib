package main

import (
	"flag"
	"fmt"
	"os"
)

func init() {

}

func switchTrue(num int) {
	switch {
	case num <= 0:
		fmt.Println("<= 0")
	case num <= 3:
		fmt.Println("<= 3")
	default:
		fmt.Println("> 3")
	}
}

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

	switchTrue(-1)
	switchTrue(1)
	switchTrue(10)
}
