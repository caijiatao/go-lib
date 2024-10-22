package exception

import (
	"github.com/pkg/errors"
)

func foo() error {
	return errors.New("something went wrong")
}

func bar() error {
	return foo() // 将堆栈信息附加到错误上
}

type DataSourceConfig struct{}

type CopyDataJob struct {
	source      *DataSourceConfig
	destination *DataSourceConfig

	err error
}

func (job *CopyDataJob) newSrc() {
	if job.err != nil {
		return
	}

	if job.source == nil {
		job.err = errors.New("source is nil")
		return
	}

	// 实例化读取数据源
}

func (job *CopyDataJob) newDst() {
	if job.err != nil {
		return
	}

	if job.destination == nil {
		job.err = errors.New("destination is nil")
		return
	}

	// 实例化写入数据源
}

func (job *CopyDataJob) copy() {
	if job.err != nil {
		return
	}

	// 复制数据 ...
}

func (job *CopyDataJob) Run() error {
	job.newSrc()
	job.newDst()

	job.copy()

	return job.err
}
