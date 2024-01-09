package inner

import (
	"fmt"
	_ "unsafe"
)

//go:linkname start outer.start
func start() {
	fmt.Println("inner start")
}
