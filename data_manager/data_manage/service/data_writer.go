package service

import (
	"context"
	"fmt"
)

type DataWriter struct {
	IDataWriter
	DataChan chan Data
}

// Run
//
//	@Description: 启动写入数据循环
func (writer *DataWriter) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-writer.DataChan:
				if !ok {
					return
				}
				err := writer.Write(data)
				// TODO log err
				fmt.Println(err)
			}
		}
	}()
}

func NewDataWriter() *DataWriter {
	return &DataWriter{
		DataChan: make(chan Data, 100),
	}
}

type IDataWriter interface {
	Write(data Data) (err error)
}

type WriterConstructor func() IDataWriter
