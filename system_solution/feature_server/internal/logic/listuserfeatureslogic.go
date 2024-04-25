package logic

import (
	"context"

	"feature_server/feature"
	"feature_server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserFeaturesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserFeaturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserFeaturesLogic {
	return &ListUserFeaturesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUserFeaturesLogic) ListUserFeatures(in *feature.UserFeatureRequest, stream feature.FeatureServer_ListUserFeaturesServer) error {
	for i := in.GetOffset(); i < 100; i++ {
		err := stream.Send(&feature.UserFeature{
			UserId: int64(i),
			Name:   "test",
		})
		if err != nil {
			return err
		}
	}
	return nil
}
