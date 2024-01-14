package outer

import (
	_ "golib/concurrency/golink/inner"
	_ "unsafe"
)

func start()

func Start() {
	start()
}
