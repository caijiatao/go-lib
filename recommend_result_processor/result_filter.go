package recommend_result_processor

import "context"

type ResultFilter interface {
	Filter(ctx context.Context, results *RecommendResults) (filterResults *RecommendResults, err error)
}

type AlreadyPurchasedFilter struct{}

type OutOfStockFilter struct{}
