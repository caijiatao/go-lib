package slice_util

type Delta struct {
	Type  string
	Value interface{}
}

type Deltas []Delta

type DeltasFIFO struct {
	Items map[string]Deltas
}

func NewDeltasFIFO() *DeltasFIFO {
	return &DeltasFIFO{
		Items: make(map[string]Deltas),
	}
}

func (f *DeltasFIFO) Add(key string, delta Delta) {
	oldDeltas := f.Items[key]
	newDeltas := append(oldDeltas, delta)
	f.Items[key] = newDeltas
}

func (f *DeltasFIFO) Get(key string) Deltas {
	return f.Items[key]
}
