package logic

import (
	"context"

	"feature_server/feature"
	"feature_server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBehaviorFeaturesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBehaviorFeaturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBehaviorFeaturesLogic {
	return &ListBehaviorFeaturesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListBehaviorFeaturesLogic) ListBehaviorFeatures(in *feature.BehaviorFeatureRequest, stream feature.FeatureServer_ListBehaviorFeaturesServer) error {
	// todo: add your logic here and delete this line

	return nil
}
