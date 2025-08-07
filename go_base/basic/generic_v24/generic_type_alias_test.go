package v24_test

import (
	"fmt"
	v24 "go/base/basic/generic_v24"
	"testing"
)

func TestSet(t *testing.T) {
	set := v24.NewSet[int](10)
	set.Add(2)
	set.Add(4)
	set.Add(6)

	if set.Len() != 3 {
		t.Fail()
	}
	if !set.Exists(2) {
		t.Fail()
	}
	set.Remove(2)
	if set.Exists(2) {
		t.Fail()
	}
	if set.Len() != 2 {
		t.Fail()
	}

	set.Range(func(a int) {
		fmt.Printf("%d\n", a)
	})
}

func BenchmarkSet(b *testing.B) {
	set := v24.NewSet[int](10)
	set.Add(2)
	set.Add(4)
	set.Add(6)
	b.ResetTimer()

	// 之前的写法
	// for i := 0; i < b.N; i++ {
	// 	set.Exists(2)
	// }

	// V1.24的写法
	for b.Loop() {
		set.Exists(2)
	}
}

// go test -v ./ -run=^TestSet$ -count=1
// go test ./ -bench=^BenchmarkSet$ -run=^$
