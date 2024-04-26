package logic

import (
	"context"

	"feature_server/feature"
	"feature_server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSubmittedJournalFeaturesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSubmittedJournalFeaturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSubmittedJournalFeaturesLogic {
	return &ListSubmittedJournalFeaturesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListSubmittedJournalFeaturesLogic) ListSubmittedJournalFeatures(in *feature.SubmittedJournalFeatureRequest, stream feature.FeatureServer_ListSubmittedJournalFeaturesServer) error {
	// todo: add your logic here and delete this line

	return nil
}
