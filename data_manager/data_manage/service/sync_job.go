package service

import (
	"airec_server/pkg/etcd_helper"
	"airec_server/pkg/goasync"
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type SyncJob struct {
	KeyPrefix     string
	mappingManage *MappingManage
}

func NewSyncJob(ctx context.Context, keyPrefix string) *SyncJob {
	job := &SyncJob{KeyPrefix: keyPrefix}
	job.startWatch(ctx)
	return job
}

type syncDataJobValue struct {
	DatasourceId uint64 `json:"datasource_id"`
}

func (j *syncDataJobValue) String() string {
	bytes, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (sj *SyncJob) jobKey(datasourceId uint64) string {
	return fmt.Sprintf("%s/%d", sj.KeyPrefix, datasourceId)
}

func (sj *SyncJob) encodeJobValue(datasourceId uint64) string {
	jv := &syncDataJobValue{DatasourceId: datasourceId}
	return jv.String()
}

func (sj *SyncJob) decodeJobValue(value []byte) syncDataJobValue {
	jv := syncDataJobValue{}
	_ = json.Unmarshal(value, &jv)
	return jv
}

func (sj *SyncJob) CreateInitSyncJob(ctx context.Context, datasourceId uint64) (err error) {
	// TODO check datasource is exists

	// TODO 检查是否已经同步完成

	// 推送一个同步数据源的任务
	_, err = etcd_helper.Put(ctx, sj.jobKey(datasourceId), sj.encodeJobValue(datasourceId))
	if err != nil {
		return err
	}
	return nil
}

func (sj *SyncJob) startWatch(ctx context.Context) {
	go func() {
		defer func() {
			r := recover()
			_ = goasync.PanicErrHandler(r)
		}()
		sj.watchInitSyncJob(ctx)
	}()
}

func (sj *SyncJob) watchInitSyncJob(ctx context.Context) {
	wch := etcd_helper.Watch(ctx, sj.KeyPrefix)
	for wr := range wch {
		for _, ev := range wr.Events {
			if ev.Type == mvccpb.PUT {
				jobVal := sj.decodeJobValue(ev.Kv.Value)
				sj.mappingManage.RunSyncDataMapping(ctx, jobVal)
			}
		}
	}
}
