package recommend_result_processor

import "context"

type ResultSort interface {
	Sort(ctx context.Context, results *RecommendResults) (sortResults *RecommendResults, err error)
}

type ResultReSort interface {
	ResultSort
}
