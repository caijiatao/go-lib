package repository

import (
	"golib/libs/test_libs"
)

func NewSegmentAPI() SegmentAPI {
	// debug 模式则直接返回一个Mock的 API模拟数据库
	if test_libs.IsDebugMode() {
		return newMockIdSegmentImpl()
	}
	return newIdSegmentRepository()
}
