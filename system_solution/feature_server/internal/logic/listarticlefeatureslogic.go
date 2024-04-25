package logic

import (
	"context"

	"feature_server/feature"
	"feature_server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticleFeaturesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListArticleFeaturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticleFeaturesLogic {
	return &ListArticleFeaturesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListArticleFeaturesLogic) ListArticleFeatures(in *feature.ArticleFeatureRequest, stream feature.FeatureServer_ListArticleFeaturesServer) error {
	// todo: add your logic here and delete this line

	return nil
}
