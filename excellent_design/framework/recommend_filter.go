package framework

type Item struct{}

type Filter interface {
	DoFilter(items []Item) []Item
}

// ConstructorFilters
//
//	@Description: 这里每次都构造出一条新的Filter，如果缓存有变化则进行更新，后续新的任务则拿到新的过滤链
//	@return []Filter
func ConstructorFilters() []Filter {
	// 这里的filters 策略可以直接从配置文件中读取，然后进行初始化
	return []Filter{
		&BlackFilter{}, // 这里内部逻辑如果有变化可以通过构造函数来实现，每次构造出来不同的逻辑
		&AlreadyBuyFilter{},
	}
}

func RunFilters(items []Item, fs []Filter) []Item {
	for _, f := range fs {
		items = f.DoFilter(items)
	}
	return items
}

type BlackFilter struct{}

func (f *BlackFilter) DoFilter(items []Item) []Item {
	// do something
	return items
}

type AlreadyBuyFilter struct{}

func (f *AlreadyBuyFilter) DoFilter(items []Item) []Item {
	// do something
	return items
}
