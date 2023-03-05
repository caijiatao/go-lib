package struct_utils

import "testing"

type aStructCopy struct {
	Name string
	Male string
}

func newAStructCopyFromAStruct(a *aStruct) *aStructCopy {
	return &aStructCopy{
		Name: a.Name,
		Male: a.Male,
	}
}

func BenchmarkCopyIntersectionStruct(b *testing.B) {
	a := &aStruct{
		Name: "test",
		Male: "test",
	}
	for i := 0; i < b.N; i++ {
		var ac aStructCopy
		CopyIntersectionStruct(a, &ac)
	}
}

func BenchmarkNormalCopyIntersectionStruct(b *testing.B) {
	a := &aStruct{
		Name: "test",
		Male: "test",
	}
	for i := 0; i < b.N; i++ {
		newAStructCopyFromAStruct(a)
	}
}
