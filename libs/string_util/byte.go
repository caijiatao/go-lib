package string_util

import "unsafe"

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func String2Bytes(s string) []byte {
	/**
	作者：高德技术
	链接：https://juejin.cn/post/7002434109247062046
	来源：稀土掘金
	著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
	*/
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
