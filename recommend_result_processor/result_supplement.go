package recommend_result_processor

import "context"

type ResultsSupplement interface {
	SupplementData(ctx context.Context, results *RecommendResults) (supplementResults *RecommendResults, err error)
}
