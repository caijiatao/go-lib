package inner

import (
	"fmt"
	_ "unsafe"
)

//go:linkname start golib/concurrency/golink/outer.start
func start() {
	fmt.Println("inner start")
}
