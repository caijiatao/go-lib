package logic

import (
	"context"

	"feature_server/feature"
	"feature_server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCooperatorFeaturesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCooperatorFeaturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCooperatorFeaturesLogic {
	return &ListCooperatorFeaturesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCooperatorFeaturesLogic) ListCooperatorFeatures(in *feature.CooperatorFeatureRequest, stream feature.FeatureServer_ListCooperatorFeaturesServer) error {
	// todo: add your logic here and delete this line

	return nil
}
