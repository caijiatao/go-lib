package runtime

import "fmt"

var (
	ReallyCrash = true
)

// 全局默认的Panic处理
var PanicHandlers = []func(interface{}){logPanic}

// 允许外部传入额外的异常处理
func HandleCrash(additionalHandlers ...func(interface{})) {
	if r := recover(); r != nil {
		for _, fn := range PanicHandlers {
			fn(r)
		}
		for _, fn := range additionalHandlers {
			fn(r)
		}
		if ReallyCrash {
			panic(r)
		}
	}
}

func logPanic(args interface{}) {
	fmt.Println("log panic", args)
}
