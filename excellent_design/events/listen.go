package events

type listener struct {
	eventObjs chan eventObj
}

// watch
//
//	@Description: 监听需要处理的内容
func (l *listener) watch() chan eventObj {
	return l.eventObjs
}

// 事件对象，可以自己定义想传递的内容
type eventObj struct{}

var (
	listeners = make([]*listener, 0)
)

func distribute(obj eventObj) {
	for _, l := range listeners {
		// 这里直接将事件对象分发出去
		l.eventObjs <- obj
	}
}
