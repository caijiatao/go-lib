package recommend_result_processor

type RecommendResultProcessor struct {
	// (AFilter && BFilter) || CFilter
	filters     []ResultFilter
	sorts       []ResultSort
	supplements []ResultsSupplement
	resort      []ResultReSort
}

func (processor *RecommendResultProcessor) Process() {
}
