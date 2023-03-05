package repository

import "golib/common"

func NewSegmentAPI() SegmentAPI {
	// debug 模式则直接返回一个Mock的 API模拟数据库
	if common.IsDebugMode() {
		return newMockIdSegmentImpl()
	}
	return newIdSegmentRepository()
}
