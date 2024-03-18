package outer

import (
	_ "golib/libs/concurrency/golink/inner"
	_ "unsafe"
)

func start()

func Start() {
	start()
}
