package define

import "errors"

var (
	RespGenerateIdErr          = errors.New("generate id err")
	RespSegmentAllocNotInitErr = errors.New("segment alloc not init")
)
