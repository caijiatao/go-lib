package define

type Segment struct {
	Cursor    uint64
	MaxCursor uint64
	Ids       []uint64
}

func NewSegment(step uint64) *Segment {
	return &Segment{
		MaxCursor: step,
		// 每次都从0开始
		Cursor: 0,
		Ids:    make([]uint64, 0, step),
	}
}
