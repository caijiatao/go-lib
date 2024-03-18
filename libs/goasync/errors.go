package goasync

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
)

func PanicErrHandler(r any) (err error) {
	switch r.(type) {
	case runtime.Error:
		err = r.(runtime.Error)
		return
	case error:
		err = r.(error)
	}
	// print stack
	if err != nil {
		log.Println(fmt.Sprintf("err: %#v,stack:%s", err, debug.Stack()))
		return err
	}
	return nil
}
