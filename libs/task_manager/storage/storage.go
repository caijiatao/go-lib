package storage

type Storage interface {
	//
	// Watch
	//  @Description: 监听任务状态
	//
	Watch()

	//
	// Put
	//  @Description: 新增任务
	//
	Put()
}
