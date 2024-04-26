package logic

import (
	"context"

	"feature_server/feature"
	"feature_server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPublishedPaperFeaturesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPublishedPaperFeaturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPublishedPaperFeaturesLogic {
	return &ListPublishedPaperFeaturesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPublishedPaperFeaturesLogic) ListPublishedPaperFeatures(in *feature.PublishedPaperFeatureRequest, stream feature.FeatureServer_ListPublishedPaperFeaturesServer) error {
	// todo: add your logic here and delete this line

	return nil
}
